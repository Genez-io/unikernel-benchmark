import asyncio, subprocess, sys
from qemu.qmp import QMPClient


class QEMUApiClient(object):
    def __init__(self, socket_path: str) -> None:
        self.socket_path = socket_path

    async def __connect__(self) -> None:
        qmp = QMPClient("osv-qemu")
        await qmp.connect(self.socket_path)
        await qmp.disconnect()

    def start_vm(self) -> None:
        asyncio.run(self.__connect__())
