FROM clearlinux:base

RUN swupd bundle-add go-basic

RUN useradd -m user
USER user

WORKDIR /home/user/
COPY backend/ backend/

WORKDIR /home/user/backend
RUN go install