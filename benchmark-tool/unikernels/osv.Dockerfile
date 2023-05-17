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

RUN mkdir /osv/.firecracker/ && \   
    wget https://github.com/firecracker-microvm/firecracker/releases/download/v0.23.0/firecracker-v0.23.0-x86_64 -O /osv/.firecracker/firecracker-x86_64 && \
    chmod a+x /osv/.firecracker/firecracker-x86_64

RUN echo "imageSizeBytes=$(wc -c ./build/release/usr.img | cut -d" " -f1)" >> /static_metrics

RUN pip install requests-unixsocket

COPY /unikernels/boot_docker_unikernel.sh /scripts/boot_docker_unikernel.sh
COPY /unikernels/osv.py /scripts/osv.py
COPY /unikernels/utils /scripts/utils

CMD ["/bin/bash", "-c", "/scripts/boot_docker_unikernel.sh \
    172.17.0.2 \
    172.16.0.2 \
    25565 \
    \"python3 /scripts/osv.py\" \
    "]
