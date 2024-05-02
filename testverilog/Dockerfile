FROM python:3.9-buster

WORKDIR /work

RUN pip install circuitgraph

RUN git clone https://github.com/circuitgraph/logiclocking.git

RUN pip3 install -e ./logiclocking
