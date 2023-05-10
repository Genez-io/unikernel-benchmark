from requests import Response
import requests_unixsocket as req
from typing import List
import json, os, signal, subprocess

class FirecrackerApiClient(object):
    def __init__(self, socket_path: str) -> None:
        self.session = req.Session()
        self.socket_path = f'http+unix://{socket_path.replace("/", "%2F")}'


    def make_put_call(self, endpoint: str, request_body: dict) -> Response:
        res = self.session.put(self.socket_path + endpoint, data=json.dumps(request_body))
        if res.status_code != 204:
            print(f'Error: {res.text}')
        return res


class FirecrackerGuestMemoryMonitor(object):
    X86_MEMORY_GAP_START = 3407872

    def __init__(self, firecracker_process, guest_mem_mib, update_interval) -> None:
        self._rss_list = []
        self._process = firecracker_process
        self._update_interval = update_interval
        self._guest_mem_mib = guest_mem_mib
        self._guest_mem_range1 = None
        self._guest_mem_range2 = None


    def reset_ranges(self):
        self._guest_mem_range1 = None
        self._guest_mem_range2 = None


    def _update_guest_memory_regions(self, address: int, size_kib: int) -> None:
        # If x86_64 guest memory exceeds 3328M, it will be split
        # in 2 regions: 3328M and the rest. We have 3 cases here
        # to recognise a guest memory region:
        #  - its size matches the guest memory exactly
        #  - its size is 3328M
        #  - its size is guest memory minus 3328M.
        if size_kib in (
            self._guest_mem_mib * 1024,
            self.X86_MEMORY_GAP_START,
            self._guest_mem_mib * 1024 - self.X86_MEMORY_GAP_START,
        ):
            if not self._guest_mem_range1:
                self._guest_mem_range1 = (address, address + size_kib * 1024)
                return True
            if not self._guest_mem_range2:
                self._guest_mem_range2 = (address, address + size_kib * 1024)
                return True
        return False


    def get_rss(self) -> List[int]:
        print(self._rss_list)
        return self._rss_list


    def run(self) -> None:
        while True:
            try:
                self._process.wait(self._update_interval)
                break
            except KeyboardInterrupt:
                os.kill(self._process.pid, signal.SIGINT)
            except subprocess.TimeoutExpired:
                # Check guest memory usage
                proc = subprocess.run(
                    ['pmap', '-xq', str(self._process.pid)],
                    universal_newlines=True,
                    stdout=subprocess.PIPE
                )
                memory = 0
                self.reset_ranges()
                for line in proc.stdout.splitlines():
                    tokens = line.split()
                    if not tokens:
                        break
                    try:
                        address = int(tokens[0].lstrip('0'), 16)
                        total_size = int(tokens[1])
                        rss = int(tokens[2])
                    except ValueError:
                        # This line doesn't contain memory related information.
                        continue

                    if self._update_guest_memory_regions(address, total_size):
                        memory += rss

                self._rss_list.append(memory)
