import subprocess, sys, time

from utils.qemu import QEMUApiClient
from utils.memory import MemoryMonitor

pass
from utils.benchmark import BenchmarkChannel


QMP_SOCKET = "qmp.sock"
api = QEMUApiClient(QMP_SOCKET)
benchmark_channel = BenchmarkChannel(25565)

# Set up the tap interface
# subprocess.run(
#     [
#         "/osv/scripts/setup_fc_networking.sh",
#         "natted",
#         "fc_tap0",
#         "172.16.0.1",
#         "qemu_tap0",
#     ]
# )

# subprocess.run(["/osv/scripts/setup-external-bridge.sh", "fc_tap0", "virbr0"])

# Send static metrics
benchmark_channel.send_static_metrics("/static_metrics")

# Start the VM
run_py = subprocess.Popen(
    [
        "./scripts/run.py",
        "--forward",
        "udp::25565-:25565",
        "--execute=/benchmark_executable",
        "-p",
        "kvm",
        "--pass-args",
        f"-qmp unix:{QMP_SOCKET},server=on,wait=on",
    ],
    stdout=sys.stdout,
    stderr=subprocess.STDOUT,
)

time.sleep(0.3)
run_sh_children = subprocess.run(
    ["ps", "-o", "pid", "--ppid", str(run_py.pid), "--noheaders"],
    text=True,
    stdout=subprocess.PIPE,
).stdout.splitlines()
qemu_pid = int(run_sh_children[0].strip())

api.start_vm()
# Mark the start of the booting process
benchmark_channel.mark_booting_start()

monitor = MemoryMonitor(run_py, qemu_pid, 1)
monitor.run()

benchmark_channel.mark_execution_end()

benchmark_channel.send_runtime_metrics(max(monitor.get_rss()) / 1024)
