FROM ubuntu:latest

RUN apt-get -y update && apt-get -y install git python3 make gcc net-tools sudo pip

RUN git clone --recurse-submodules https://github.com/cloudius-systems/osv /osv

WORKDIR /osv

RUN python3 ./scripts/setup.py

RUN ./scripts/build

COPY benchmark-executable /benchmark-executable
COPY benchmark-framework /benchmark-framework

RUN make -C /benchmark-executable

RUN ./scripts/manifest_from_host.sh -w ../benchmark-executable/benchmark_executable && ./scripts/build --append-manifest

RUN echo "imageSizeBytes=$(wc -c ./build/release/usr.img | cut -d" " -f1)" >> /static_metrics

RUN pip install qemu.qmp asyncio

COPY /unikernels/osv/qemu/osv.py /scripts/osv.py
COPY /unikernels/utils /scripts/utils

CMD ["/bin/bash", "-c", "python3 /scripts/osv.py"]
