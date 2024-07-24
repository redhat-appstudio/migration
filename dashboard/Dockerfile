FROM registry.access.redhat.com/ubi9/go-toolset:1.18.10-4 as builder

COPY ./ ./
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download
RUN go build -o bin/generator .

FROM registry.access.redhat.com/ubi9-minimal:9.4-1194

COPY --from=builder /opt/app-root/src/bin/generator /
ENTRYPOINT [ "/generator" ]
