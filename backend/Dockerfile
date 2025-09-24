FROM golang:1.25

WORKDIR /usr/src/app

# install ffmpeg + ffprobe via apt
RUN apt-get update && \
    apt-get install -y ffprobe && \
    rm -rf /var/lib/apt/lists/*

# Go modules
COPY go.mod go.sum ./
RUN go mod download

# App source
COPY . .

# Build
RUN go build -v -o /usr/local/bin/app ./main.go

# Expose and run
EXPOSE 8080
CMD ["app"]

