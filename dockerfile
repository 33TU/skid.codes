FROM fedora:latest

RUN dnf update -y
RUN dnf install -y golang

RUN useradd -m user
USER user

WORKDIR /home/user/
COPY backend/ backend/

WORKDIR /home/user/backend
RUN go install