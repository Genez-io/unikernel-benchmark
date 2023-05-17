import subprocess, sys, os, signal

from utils.benchmark import BenchmarkChannel

benchmark_channel = BenchmarkChannel(25565)

benchmark_channel.send_static_metrics('/static_metrics')

qemu = subprocess.Popen(["./run.sh", "-n", "-k",
                         "/elfloader/apps/app-elfloader/build/app-elfloader_qemu-x86_64",
                         "-r", "/dynamic-apps/benchmark-executable/", "/bin/benchmark_executable"],
                         stdout=sys.stdout, stderr=subprocess.STDOUT)

# Mark the start of the booting process
benchmark_channel.mark_booting_start()

try:
    qemu.wait(10)
except KeyboardInterrupt:
    os.kill(qemu.pid, signal.SIGINT)
except subprocess.TimeoutExpired:
    os.kill(qemu.pid, signal.SIGINT)

# Mark the end of the execution
benchmark_channel.mark_execution_end()

benchmark_channel.send_runtime_metrics(0.0)