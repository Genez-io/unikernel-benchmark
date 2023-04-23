import requests_unixsocket as req
import json, subprocess, sys, os, signal, time, socket

class FirecrackerApiClient(object):
    def __init__(self, socket_path: str):
        self.session = req.Session()
        self.socket_path = f'http+unix://{socket_path.replace("/", "%2F")}'

    def make_put_call(self, endpoint: str, request_body: dict):
        res = self.session.put(self.socket_path + endpoint, data=json.dumps(request_body))
        if res.status_code != 204:
            print(f'Error: {res.text}')
        return res.status_code
    
class BenchmarkChannel(object):
    def __init__(self, port: int):
        self.conn = None
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
            s.bind(("", port))
            s.listen()
            conn, addr = s.accept()
            self.conn = conn

    def send_static_metrics(self, metrics_file: str):
        metrics = {}
        # Read contents from metrics file and send it to the client
        with open(metrics_file, "r") as f:
            for line in f:
                key, value = line.strip().split("=")
                try:
                    value = int(value)
                except ValueError:
                    pass

                metrics[key] = value

        self.conn.sendall(json.dumps(metrics).encode("utf-8"))

    def mark_booting_start(self):
        self.conn.close()
            


SOCKET_PATH = '/osv/.firecracker/api.socket'
FIRECRACKER_PATH = '/osv/.firecracker/firecracker-x86_64'

api = FirecrackerApiClient(SOCKET_PATH)
benchmark_channel = BenchmarkChannel(25565)

# Set up the tap interface
subprocess.run(['/osv/scripts/setup_fc_networking.sh', 'natted', 'fc_tap0', '172.16.0.1'])
# Convert osv image to raw
subprocess.run(['qemu-img', 'convert', '-O', 'raw', '/osv/build/last/usr.img', '/osv/build/last/usr.raw'])

# Start firecracker
fc = subprocess.Popen([FIRECRACKER_PATH, '--api-sock', SOCKET_PATH],
                  stdout=sys.stdout, stderr=subprocess.STDOUT)
while not os.path.exists(SOCKET_PATH):
    time.sleep(0.01)

# Set boot source
api.make_put_call('/boot-source', {
       'kernel_image_path': '/osv/build/last/loader-stripped.elf',
       'boot_args': '--ip=eth0,172.16.0.2,255.255.255.252 --defaultgw=172.16.0.1 --nameserver=172.16.0.1 --nopci /benchmark_executable'
    })

# Set disk
api.make_put_call('/drives/rootfs', {
        'drive_id': 'rootfs',
        'path_on_host': '/osv/build/last/usr.raw',
        'is_root_device': False,
        'is_read_only': False
    })

# Set network interface
api.make_put_call('/network-interfaces/eth0', {
        'iface_id': 'eth0',
        'host_dev_name': 'fc_tap0',
        'guest_mac': '52:54:00:12:34:56',
        'rx_rate_limiter': {
            'bandwidth': {
            'size': 0,
            'refill_time': 0
            },
            'ops': {
            'size': 0,
            'refill_time': 0
            }
        },
        'tx_rate_limiter': {
            'bandwidth': {
            'size': 0,
            'refill_time': 0
            },
            'ops': {
            'size': 0,
            'refill_time': 0
            }
        }
    })

# Set machine configuration
api.make_put_call('/machine-config', {
        'vcpu_count': 1,
        'mem_size_mib': 1024,
        'ht_enabled': False
    })

benchmark_channel.send_static_metrics('/static_metrics')

# Start the booting process
api.make_put_call('/actions', {'action_type': 'InstanceStart'})

# Mark the start of the booting process
benchmark_channel.mark_booting_start()

try:
    fc.wait()
except KeyboardInterrupt:
    os.kill(fc.pid, signal.SIGINT)
