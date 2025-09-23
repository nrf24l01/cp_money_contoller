from dotenv import load_dotenv
from os import getenv

load_dotenv()

# Server config
WORKER_SERVER = getenv("WORKER_SERVER", "127.0.0.1")
WORKER_TASK_QUEUE = getenv("WORKER_TASK_QUEUE", "task_queue")
WORKER_SUBMIT_HOST = getenv("WORKER_SUBMIT_HOST", "http://127.0.0.1:1327")

# IMAP config
IMAP_HOST = getenv("IMAP_HOST")
IMAP_USER = getenv("IMAP_USER")
IMAP_PASSWORD = getenv("IMAP_PASSWORD")