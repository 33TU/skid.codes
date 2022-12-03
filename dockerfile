FROM clearlinux:base

# Install depencies
RUN swupd bundle-add go-basic

# Add and switch user
RUN useradd -m user
#USER user

# Install go modules
COPY backend/ /home/user/backend/
WORKDIR /home/user/backend
RUN go install