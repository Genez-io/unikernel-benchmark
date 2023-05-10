import json, socket
from typing import Any

class TypedDict(dict):
    def __setitem__(self, __key: Any, __value: Any) -> None:
        model = self.__dict__

        if __key not in model:
            raise KeyError(f'Key {__key} does not exist in model')
        
        if not isinstance(__value, type(model[__key])):
            raise ValueError(f'Value {__value} is not of type {type(model[__key])}')
        
        super().__setitem__(__key, __value)


class StaticMetrics(TypedDict):
    def __init__(self) -> None:
        self.imageSizeBytes: int = int(0)

class RuntimeMetrics(TypedDict):
    def __init__(self) -> None:
        self.totalMemoryUsageMiB: float = float(0)

class BenchmarkChannel(object):
    def __init__(self, port: int) -> None:
        self.conn = None
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
            s.bind(('', port))
            s.listen()
            conn, addr = s.accept()
            # Disable Nagle's algorithm
            conn.setsockopt(socket.IPPROTO_TCP, socket.TCP_NODELAY, 1)

            self.conn = conn


    def __del__(self) -> None:
        self.conn.close()


    def send_tcp_message(self, data: bytes) -> None:
        data_len = len(data)
        self.conn.sendall(data_len.to_bytes(4, byteorder='little'))
        self.conn.sendall(data)


    def send_static_metrics(self, metrics_file: str) -> None:
        metrics = StaticMetrics()
        # Read contents from metrics file and send it to the client
        with open(metrics_file, 'r') as f:
            for line in f:
                key, value = line.strip().split('=')
                try:
                    value = int(value)
                except ValueError:
                    pass

                try:
                    metrics[key] = value
                except KeyError:
                    print(f'Field {key} is not a valid static metric')
                except ValueError:
                    print(f'Value {value} is not of type {type(metrics[key])}')

        self.send_tcp_message(json.dumps(metrics).encode('utf-8'))


    def mark_booting_start(self):
        self.send_tcp_message('start_booting'.encode('utf-8'))


    def mark_execution_end(self):
        self.send_tcp_message('execution_ended'.encode('utf-8'))


    def send_runtime_metrics(self, max_rss_mib: int) -> None:
        metrics = RuntimeMetrics()
        metrics['totalMemoryUsageMiB'] = max_rss_mib

        self.send_tcp_message(json.dumps(metrics).encode('utf-8'))