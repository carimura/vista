FROM golang:latest as build-stage
WORKDIR /src
RUN go get github.com/minio/minio-go
ADD . /src
RUN go build -o func


FROM carimura/ffmpeg
WORKDIR /data
COPY --from=build-stage /src/* /data/

ENTRYPOINT ["./func"]
