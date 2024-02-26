# FROM golang:1.22-bullseye
# FROM gocv/opencv:4.8.1-ubuntu-22.04
FROM gocv/opencv:4.8.1
# FROM arm32v7/golang:1.22-alpine

# RUN apk update && apk add make gcc g++ pkgconfig opencv-dev g++-arm-linux-eabi gcc-arm-linux-eabi binutils-arm-linux-eabi
RUN apt-get update && apt-get install -y build-essential pkg-config
# gcc-arm-linux-gnueabihf g++-arm-linux-gnueabihf
# RUN apt-get update && apt-get install -y build-essential pkg-config gcc-10-arm-linux-gnueabihf g++-10-arm-linux-gnueabihf
# libopencv-dev 

# RUN wget https://github.com/prepkg/opencv-raspberrypi/releases/download/4.8.0/opencv.deb
# RUN dpkg -i opencv.deb

# RUN apk update && apk add build-base pkgconfig opencv

# RUN wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
# RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
# RUN export PATH=$PATH:/usr/local/go/bin

# RUN wget -P ~ https://musl.cc/aarch64-linux-musl-cross.tgz
RUN wget -P ~ https://musl.cc/armv7l-linux-musleabihf-cross.tgz
# RUN tar -xvf ~/aarch64-linux-musl-cross.tgz -C /usr/local
RUN tar -xvf ~/armv7l-linux-musleabihf-cross.tgz -C /usr/local
# RUN ln -f -s /usr/local/aarch64-linux-musl/libc.so /usr/local/aarch64-linux-musl-cross/aarch64-linux-musl/lib/ld-musl-aarch64.so.1
RUN ln -f -s /usr/local/armv7l-linux-musleabihf-cross/libc.so /usr/local/armv7l-linux-musleabihf-cross/armv7l-linux-musleabihf/lib/ld-musl-armhf.so.1
# RUN export PATH=$PATH:/usr/local/aarch64-linux-musl-cross/bin
RUN export PATH=$PATH:/usr/local/armv7l-linux-musleabihf-cross/bin

# CMD ["ls", "-la", "/usr/local/armv7l-linux-musleabihf-cross/bin/"]

ENV CGO_ENABLED=1
ENV GOARCH=arm
ENV GOARM=7
ENV GOOS=linux
ENV CC=/usr/local/armv7l-linux-musleabihf-cross/bin/armv7l-linux-musleabihf-gcc
ENV CXX=/usr/local/armv7l-linux-musleabihf-cross/bin/armv7l-linux-musleabihf-g++
# ENV CC=arm-linux-gnueabihf-gcc-10
# ENV CXX=arm-linux-gnueabihf-g++-10

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN go mod download

# RUN go install github.com/cosmtrek/air@latest

# CMD ["which", "arm-none-eabi-gcc"]
# CMD ["air", "-c", "build/air.toml"]
CMD ["go", "build", "-v", "-trimpath", "-ldflags", "-extldflags -static"]
# CMD ["lsb_release"]
# CMD ["pkg-config --cflags --libs opencv"]
# CMD ["pkg-config", "--cflags", "--libs", "opencv"]
