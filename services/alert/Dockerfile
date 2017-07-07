FROM golang:latest as build-stage
WORKDIR /src
RUN go get github.com/ChimeraCoder/anaconda
ADD . /src
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o func .

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /func
COPY --from=build-stage /src/* /func/
ENTRYPOINT ["./func"]

