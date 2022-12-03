FROM clearlinux:base

# Install depencies
RUN swupd bundle-add go-basic

# Copy backend to dest
COPY backend/ /run/backend/
WORKDIR /run/backend

# Install go modules
RUN go install