import os
import sys
# Ensure generated protobuf modules can be imported
sys.path.insert(0, os.path.dirname(__file__))
# Dynamically load generated protobuf modules to ensure correct resolution
import importlib.util
base_dir = os.path.dirname(__file__)
# Dynamically load Worker_pb2
spec_pb2 = importlib.util.spec_from_file_location(
    'Worker_pb2', os.path.join(base_dir, 'Worker_pb2.py')
)
assert spec_pb2 and spec_pb2.loader
Worker_pb2 = importlib.util.module_from_spec(spec_pb2)
spec_pb2.loader.exec_module(Worker_pb2)
# Dynamically load Worker_pb2_grpc
spec_pb2_grpc = importlib.util.spec_from_file_location(
    'Worker_pb2_grpc', os.path.join(base_dir, 'Worker_pb2_grpc.py')
)
assert spec_pb2_grpc and spec_pb2_grpc.loader
Worker_pb2_grpc = importlib.util.module_from_spec(spec_pb2_grpc)
spec_pb2_grpc.loader.exec_module(Worker_pb2_grpc)
# Expose message and stub classes
RegisterRequest = Worker_pb2.RegisterRequest
RegisterResponse = Worker_pb2.RegisterResponse
WorkerServiceStub = Worker_pb2_grpc.WorkerServiceStub
import grpc


class WorkerClient:
    def __init__(self, server_address: str):
        self.server_address = server_address
        self.channel = grpc.insecure_channel(server_address)
        self.stub = WorkerServiceStub(self.channel)
        
        self.uuid = None
        self.key = None

    def close(self):
        self.channel.close()
    
    def register_worker(self, worker_id: int) -> RegisterResponse:
        """
        Registers a worker by ID.
        Returns a RegisterResponse containing uuid, key, and confirmed ID.
        """
        request = RegisterRequest(id=worker_id)
        response = self.stub.RegisterWorker(request)
        self.uuid = response.uuid
        self.key = response.key
        return response

    def get_task(self):
        request = Worker_pb2.GetTaskRequest(worker_uuid=self.uuid, worker_key=self.key)
        resp = self.stub.GetTask(request)
        if resp.uuid == "":
            return None
        return resp
    
    def change_status(self, task_uuid: str, payload: str, logs: str, status: str, unix_time: int):
        request = Worker_pb2.ChangeStatusRequest(
            uuid=task_uuid,
            payload=payload,
            logs=logs,
            status=status,
            unix_time=unix_time
        )
        resp = self.stub.ChangeStatus(request)
        return resp.ok == True
    

if __name__ == '__main__':
    import os

    SERVER = os.getenv('WORKER_SERVER', 'localhost:50051')
    WORKER_ID = int(os.getenv('WORKER_ID', '1'))
    client = WorkerClient(SERVER)
    resp = client.register_worker(WORKER_ID)
    print(f"Registered worker: uuid={resp.uuid}, key={resp.key}, id={resp.id}")
    print(client.get_task())
    client.close()
