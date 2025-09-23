from json import loads
from .rabbitMQ_client import RabbitMQClient
from .logger import ThreadSafeLogger
from requests import put
from time import sleep


class TaskClient(RabbitMQClient):
    def __init__(self, host: str, queue: str, submit_host: str, logger: ThreadSafeLogger):
        super().__init__(host=host, queue=queue)
        self.logger = logger
        self.logger.info(f"Connecting to RabbitMQ at {host}:{queue}")
        self.connect()
        self.logger.info("Connected")
        self.task = None
        self.submit_host = submit_host
    
    def get_task(self):
        self.task = None
        while self.task is None:
            self.task = self.get_one()
            sleep(1)
        
        if not isinstance(self.task, dict):
            self.task = loads(self.task)
        return self.task
    
    def change_status(self, task_uuid: str, payload: dict, logs: list[str], status: str, unix_time: int):
        json_payload = {
            "result": payload,
            "logs": logs,
            "status": status,
            "unix_time": unix_time
        }
        rq = put(self.submit_host+"/worker/task/"+task_uuid, json=json_payload)

        if rq.status_code != 200:
            self.logger.error(f"Failed to update task status: {rq.status_code} - {rq.text}")
