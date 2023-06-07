import os, signal, subprocess
from typing import List


class MemoryMonitor(object):
    def __init__(self, waiting_process, pid, update_interval) -> None:
        self._rss_list = []
        self._process = waiting_process
        self._update_interval = update_interval
        self._pid = pid

    def get_rss(self) -> List[int]:
        # print(self._rss_list)
        return self._rss_list

    def run(self) -> None:
        tries_left = 10
        while True:
            try:
                self._process.wait(self._update_interval)
                break
            except KeyboardInterrupt:
                os.kill(self._pid, signal.SIGINT)
            except subprocess.TimeoutExpired:
                if tries_left <= 0:
                    os.kill(self._pid, signal.SIGINT)
                    break

                # Compute process RSS
                pmap = subprocess.run(
                    ["pmap", "-x", str(self._pid)],
                    text=True,
                    stdout=subprocess.PIPE,
                )
                line = pmap.stdout.splitlines()[-1]
                rss = int(line.split()[3])
                self._rss_list.append(rss)

                tries_left -= 1
