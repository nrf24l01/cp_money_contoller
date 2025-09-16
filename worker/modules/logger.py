import logging
import threading
from typing import List


class InMemoryHandler(logging.Handler):
    """Хендлер для хранения логов в памяти."""
    def __init__(self) -> None:
        super().__init__()
        self._records: List[str] = []
        self._lock = threading.Lock()
        self.setFormatter(logging.Formatter(
            "%(asctime)s [%(levelname)s] %(message)s"
        ))

    def emit(self, record: logging.LogRecord) -> None:
        msg = self.format(record)
        with self._lock:
            self._records.append(msg)

    def get_logs(self) -> str:
        with self._lock:
            return "\n".join(self._records)


class ThreadSafeLogger:
    def __init__(self, name: str = "MyLogger") -> None:
        self.logger = logging.getLogger(name)
        self.logger.setLevel(logging.DEBUG)

        self.memory_handler = InMemoryHandler()
        self.logger.addHandler(self.memory_handler)

        console_handler = logging.StreamHandler()
        console_handler.setFormatter(logging.Formatter(
            "[%(threadName)s] %(levelname)s: %(message)s"
        ))
        self.logger.addHandler(console_handler)

    def _log(self, level: int, *args: object) -> None:
        msg = " ".join(map(str, args))
        self.logger.log(level, msg)

    def debug(self, *args: object) -> None:
        self._log(logging.DEBUG, *args)

    def info(self, *args: object) -> None:
        self._log(logging.INFO, *args)

    def warn(self, *args: object) -> None:
        self._log(logging.WARNING, *args)

    def error(self, *args: object) -> None:
        self._log(logging.ERROR, *args)

    def critical(self, *args: object) -> None:
        self._log(logging.CRITICAL, *args)

    def get_full_log(self) -> str:
        return self.memory_handler.get_logs()
    
    def clear_log(self) -> None:
        self.memory_handler._records.clear()