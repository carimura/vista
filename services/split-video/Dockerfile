FROM golang:latest as build-stage
WORKDIR /src
ADD . /src
RUN go build -o myapp

FROM carimura/ffmpeg
WORKDIR /data
COPY --from=build-stage /src/* /data/

ENTRYPOINT ["./myapp"]
