FROM golang:1.11 as builder
RUN go get -u github.com/golang/dep/cmd/dep
WORKDIR /go/src/github.com/camptocamp/prometheus-orchestrators-sd
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only
COPY . .
RUN make posd

FROM scratch
COPY --from=builder /etc/ssl /etc/ssl
COPY --from=builder /go/src/github.com/camptocamp/prometheus-orchestrators-sd/posd /
ENTRYPOINT ["/posd"]
CMD [""]
