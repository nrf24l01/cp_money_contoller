from json import dumps, loads
from modules import WorkerClient, ImapClient, ThreadSafeLogger
from config import WORKER_SERVER, WORKER_ID, IMAP_HOST, IMAP_USER, IMAP_PASSWORD
from tasks import register_user
from time import sleep, time
from tasks import build_config
from selenium import webdriver
from threading import Thread

def wait_for_task(client: WorkerClient, logger: ThreadSafeLogger):
    logger.info("Waiting for task...")
    while True:
        task = client.get_task()
        if task:
            logger.info(f"Received task: {task}")
            return task
        logger.debug("No task available, retrying in 10 seconds...")
        sleep(10)

def update_logs(client: WorkerClient, logger: ThreadSafeLogger, task_uuid: str, res: dict):
    while True:
        sleep(3)
        logs = logger.get_full_log()
        logs = dumps(logs.split("\n"))
        client.change_status(task_uuid=task_uuid, payload=dumps(res), logs=logs, status=dumps({"status": res.get("status"), "details": res.get("details", {})}), unix_time=int(time()))
        if res.get("status") is True:
            client.change_status(task_uuid=task_uuid, payload=dumps(res), logs=logs, status=dumps({"status": res.get("status"), "details": res.get("details", {})}), unix_time=int(time()))
            return

def main():
    logger = ThreadSafeLogger("Main")
    client = WorkerClient(WORKER_SERVER)
    
    if WORKER_ID is None:
        logger.error("WORKER_ID is not set in environment variables")
        return
    
    try:
        worker_id = int(WORKER_ID)
    except ValueError:
        logger.error("WORKER_ID must be an integer")
        return

    response = client.register_worker(worker_id)
    logger.info(f"Registered worker: uuid={response.uuid}, key={response.key}, id={response.id}")
    
    logger.info("Setting up Selenium WebDriver...")
    options, service = build_config(headless=True)
    driver = webdriver.Chrome(service=service, options=options)
    logger.info("Selenium WebDriver is set up")
    
    logger.info("Starting ImapClient...")
    imap_client = ImapClient(IMAP_HOST, IMAP_USER, IMAP_PASSWORD)
    logger.info("Started ImapClient")
    
    task = wait_for_task(client, logger)
    if task:
        logger.info(f"Received task: {task}")
    else:
        logger.info("No task available")
    
    payload = loads(task.payload)
    au = payload["auth_data"]
    
    logger.clear_log()  
    
    res = {}
    lt = Thread(target=update_logs, args=(client, logger, task.uuid, res), name="LogUpdaterThread")
    lt.start()
    if lt.is_alive():
        logger.info("Logger has started")
    if task.task == "register":
        logger.info("Starting registration task...")
        t = Thread(target=register_user, args=(driver, logger, imap_client, au["invite"], au["email"], au["password"], au["name"], au["birthday"], au["learn_place"], au["grade"], res), name="WorkThread")
        t.start()
        t.join()
    else:
        logger.error(f"Unknown task type: {task.task}")
        client.change_status(task_uuid=task.uuid, payload=dumps({"error": f"Unknown task type: {task.task}"}), logs=logger.get_full_log(), status="error", unix_time=int(time()))

    logger.info("Task completed, waiting for log updater to finish...")
    lt.join()
    logger.info("Log updater finished, exiting.")

if __name__ == "__main__":
    main()