FROM nvidia/cuda:11.8.0-cudnn8-devel-ubuntu20.04

ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=Asia/Seoul

RUN apt-get update && apt-get install -y \
    python3.8 \
    python3-pip \
    python3.8-dev \
    python3.8-venv \
    git \
    wget \
    curl \
    vim \
    tzdata \
    openssh-server \
    && apt-get clean

RUN mkdir /var/run/sshd
RUN echo "PermitRootLogin yes" >> /etc/ssh/sshd_config
RUN echo "PubkeyAuthentication yes" >> /etc/ssh/sshd_config
RUN echo "PasswordAuthentication no" >> /etc/ssh/sshd_config

RUN python3.8 -m pip install --upgrade pip \
    && pip3 install --no-cache-dir tensorflow \
    && pip3 install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cu118 \
    && pip3 install jupyter \
    && mkdir /workspace

WORKDIR /workspace

EXPOSE 8888 22

CMD ["/bin/sh", "-c", "/usr/sbin/sshd && jupyter notebook --ip=0.0.0.0 --allow-root --no-browser --NotebookApp.token=''"]