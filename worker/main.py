from json import dumps, loads
from modules import TaskClient, ImapClient, ThreadSafeLogger
from config import WORKER_SERVER, WORKER_TASK_QUEUE, WORKER_SUBMIT_HOST, IMAP_HOST, IMAP_USER, IMAP_PASSWORD
from tasks import register_user
from time import sleep, time
from tasks import build_config
from selenium import webdriver
from threading import Thread
import signal

def update_logs(client: TaskClient, logger: ThreadSafeLogger, task_uuid: str, res: dict):
    while True:
        sleep(3)
        logs = logger.get_full_log()
        logs = logs.split("\n")
        client.change_status(task_uuid=task_uuid, payload=res, logs=logs, status=res.get("status", "working"), unix_time=int(time()))
        if res.get("success", None) is not None:
            client.change_status(task_uuid=task_uuid, payload=res, logs=logs, status=res.get("status", "DONE"), unix_time=int(time()))
            logger.info("Log updater thread exiting...")
            return

def main():
    logger = ThreadSafeLogger("Main")
    client = TaskClient(WORKER_SERVER, WORKER_TASK_QUEUE, WORKER_SUBMIT_HOST, logger)
    
    logger.info("Setting up Selenium WebDriver...")
    options, service = build_config(headless=False)
    driver = webdriver.Chrome(service=service, options=options)
    logger.info("Selenium WebDriver is set up")
    
    logger.info("Starting ImapClient...")
    imap_client = ImapClient(IMAP_HOST, IMAP_USER, IMAP_PASSWORD)
    logger.info("Started ImapClient")
    
    while running:
        logger.info("Waiting for task")
        task = client.get_task()
        if task:
            logger.info(f"Received task: {task}")
        else:
            logger.info("No task available")
            return

        payload = task["payload"]
        au = payload
        
        logger.clear_log()  
        
        res = {}
        lt = Thread(target=update_logs, args=(client, logger, task["uuid"], res), name="LogUpdaterThread", daemon=True)
        lt.start()
        if lt.is_alive():
            logger.info("Logger has started")
        if task["type"] == "register":
            logger.info("Starting registration task...")
            t = Thread(target=register_user, args=(driver, logger, imap_client, au["invite"], au["email"], au["password"], au["name"], au["birthday"], au["learn_place"], au["grade"], res), name="WorkThread", daemon=True)
            t.start()
            t.join()
        else:
            logger.error(f"Unknown task type: {task["type"]}")
            client.change_status(task_uuid=task["uuid"], payload={"error": f"Unknown task type: {task["type"]}"}, logs=logger.get_full_log().split("\n"), status="error", unix_time=int(time()))

        logger.info("Task completed, waiting for log updater to finish...")
        lt.join()
        logger.info("Log updater finished, exiting.")
        sleep(1)

    driver.quit()
    return

running = True

def signal_handler(signum, frame):
    global running
    running = False

signal.signal(signal.SIGTERM, signal_handler)

if __name__ == "__main__":
    main()
    exit(0)