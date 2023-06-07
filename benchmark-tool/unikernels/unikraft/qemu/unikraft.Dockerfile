FROM ubuntu:latest

RUN apt-get -y update && apt-get -y install git make gcc sudo libncurses-dev bison flex wget unzip python3 iptables qemu-system-x86 iproute2 pip

RUN mkdir elfloader elfloader/apps elfloader/libs && \
    git clone https://github.com/unikraft/app-elfloader.git /elfloader/apps/app-elfloader && \
    git clone https://github.com/unikraft/lib-libelf.git /elfloader/libs/libelf && \
    git clone https://github.com/unikraft/lib-lwip.git /elfloader/libs/lwip && \
    git clone https://github.com/unikraft/unikraft.git /elfloader/unikraft && \
    git clone https://github.com/unikraft/dynamic-apps.git /dynamic-apps && \
    git clone https://github.com/unikraft/run-app-elfloader.git /run-app-elfloader

# RUN sed '/^\timply LIBPOSIX_USER/a \\timply LIBVFSCORE_AUTOMOUNT_ROOTFS' /elfloader/apps/app-elfloader/Config.uk > Config.tmp && \
#     sed '/^\timply LIBPOSIX_USER/a \\timply PAGING' Config.tmp > /elfloader/apps/app-elfloader/Config.uk && \
#     sed '/^\timply LIBPOSIX_USER/a \\timply VIRTIO_PCI' /elfloader/apps/app-elfloader/Config.uk > Config.tmp && \
#     sed '/^\timply LIBPOSIX_USER/a \\timply PLAT_KVM' Config.tmp > /elfloader/apps/app-elfloader/Config.uk && \
#     mv Config.tmp /elfloader/apps/app-elfloader/Config.uk
#     # rm Config.tmp

WORKDIR /elfloader/apps/app-elfloader
COPY unikernels/unikraft/unikraft.config .config

RUN make

RUN echo "imageSizeBytes=$(wc -c /elfloader/apps/app-elfloader/build/app-elfloader_qemu-x86_64 | cut -d" " -f1)" >> /static_metrics

COPY benchmark-executable /benchmark-executable
COPY benchmark-framework /benchmark-framework

WORKDIR /dynamic-apps

RUN make -C /benchmark-executable && \
    mkdir /dynamic-apps/benchmark-executable /dynamic-apps/benchmark-executable/bin && \
    ./extract.sh /benchmark-executable/benchmark_executable /dynamic-apps/benchmark-executable && \
    cp /benchmark-executable/benchmark_executable /dynamic-apps/benchmark-executable/bin/benchmark_executable

WORKDIR /run-app-elfloader

RUN pip install qemu.qmp asyncio

COPY /unikernels/unikraft/qemu/unikraft.py /scripts/unikraft.py
COPY /unikernels/utils /scripts/utils

CMD ["/bin/bash", "-c", "/scripts/utils/forward_udp_to_unikernel.sh \
    172.17.0.2 \
    172.44.0.2 \
    25565 \
    \"python3 /scripts/unikraft.py\" \
    "]