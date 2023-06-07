import subprocess, sys, time
from utils.memory import MemoryMonitor

from utils.benchmark import BenchmarkChannel
from utils.qemu import QEMUApiClient

SOCKET_PATH = "/unikraft/.firecracker/api.socket"

# QMP_SOCKET = "qmp.sock"
# api = QEMUApiClient(QMP_SOCKET)
benchmark_channel = BenchmarkChannel(25565)

benchmark_channel.send_static_metrics("/static_metrics")

run_sh = subprocess.Popen(
    [
        "./run.sh",
        "-n",
        "-k",
        "/elfloader/apps/app-elfloader/build/app-elfloader_qemu-x86_64",
        "-r",
        "/dynamic-apps/benchmark-executable/",
        "/bin/benchmark_executable",
    ],
    stdout=sys.stdout,
    stderr=subprocess.STDOUT,
)
# Mark the start of the booting process
benchmark_channel.mark_booting_start()

time.sleep(0.3)
run_sh_children = subprocess.run(
    ["ps", "-o", "pid", "--ppid", str(run_sh.pid), "--noheaders"],
    text=True,
    stdout=subprocess.PIPE,
).stdout.splitlines()
qemu_pid = int(run_sh_children[0].strip())

monitor = MemoryMonitor(run_sh, qemu_pid, 1)
monitor.run()

# Mark the end of the execution
benchmark_channel.mark_execution_end()

benchmark_channel.send_runtime_metrics(max(monitor.get_rss()) / 1024)
