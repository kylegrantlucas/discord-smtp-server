FROM golang:1.15-alpine as builder
WORKDIR /build
COPY *.go go.* /build/
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -o discord-smtp .

FROM scratch
COPY --from=builder /build/discord-smtp .
ENTRYPOINT [ "./discord-smtp" ]