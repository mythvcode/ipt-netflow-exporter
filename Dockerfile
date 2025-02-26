ARG GO_VERSION

# Build userspace ebpf program
FROM golang:${GO_VERSION} as gobuilder

WORKDIR /build
ADD . ./
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o ipt-netflow-exporter ./cmd/iptnetflowexporter

# Copy builded programs to alpine image
FROM alpine:latest

LABEL description="ip-netflow prometheus exporter"
RUN addgroup --gid 39355 ipt_netflow_exporter && \
    adduser -h /app -s /bin/sh -G ipt_netflow_exporter -u 39355 -D ipt_netflow_exporter
WORKDIR /app/
COPY --from=gobuilder /build/ipt-netflow-exporter .

USER storm_control

ENTRYPOINT ["/app/ipt-netflow-exporter"]
