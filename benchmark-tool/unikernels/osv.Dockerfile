FROM ubuntu:latest

RUN apt-get -y update && apt-get -y install git python3

RUN git clone --recurse-submodules https://github.com/cloudius-systems/osv /osv

WORKDIR /osv

RUN python3 ./scripts/setup.py

RUN ./scripts/build

RUN ./scripts/build --append-manifest

RUN ./scripts/manifest_from_host.sh -w sleep && ./scripts/build --append-manifest

CMD ["/osv/scripts/run.py", "-e", "sleep 3"]