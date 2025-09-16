from dotenv import load_dotenv
from os import getenv

load_dotenv()

# Socket server conifg
WORKER_SERVER = getenv("WORKER_SERVER", "127.0.0.1:50051")
WORKER_ID = getenv("WORKER_ID")

# IMAP config
IMAP_HOST = getenv("IMAP_HOST")
IMAP_USER = getenv("IMAP_USER")
IMAP_PASSWORD = getenv("IMAP_PASSWORD")