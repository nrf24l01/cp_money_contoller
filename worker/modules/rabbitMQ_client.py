import pika
from typing import Optional


class RabbitMQClient:
    def __init__(
        self,
        host: str = "localhost",
        queue: str = "test_queue"
    ) -> None:
        self.host: str = host
        self.queue: str = queue
        self.connection: Optional[pika.BlockingConnection] = None
        self.channel: Optional[pika.adapters.blocking_connection.BlockingChannel] = None

    def connect(self) -> None:
        """Устанавливает соединение и объявляет очередь."""
        self.connection = pika.BlockingConnection(pika.ConnectionParameters(host=self.host))
        self.channel = self.connection.channel()
        self.channel.queue_declare(queue=self.queue, durable=False)

    def get_one(self) -> Optional[str]:
        """Берет одно сообщение из очереди и удаляет его."""
        if self.channel is None:
            raise RuntimeError("Сначала вызови connect()")

        method_frame, header_frame, body = self.channel.basic_get(queue=self.queue, auto_ack=True)
        if method_frame:
            return body.decode("utf-8")
        return None

    def close(self) -> None:
        """Закрывает соединение."""
        if self.connection and not self.connection.is_closed:
            self.connection.close()


if __name__ == "__main__":
    consumer = RabbitMQClient(host="localhost", queue="task_queue")
    consumer.connect()
    message = consumer.get_one()
    if message:
        print(f"[x] Получено: {message}")
    else:
        print("[*] Очередь пуста")
    consumer.close()
