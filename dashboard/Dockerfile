FROM registry.access.redhat.com/ubi9/go-toolset:1.18.10-4 as builder

COPY ./ ./
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
USER 0
RUN echo $PWD && ls -l -a && yum install -y git
RUN git config --global --add safe.directory $PWD && git log --graph --oneline --all --decorate && git fetch && git log --graph --oneline --all
RUN go mod download
RUN go build -o bin/generator .



FROM registry.access.redhat.com/ubi9-minimal:9.2-484

COPY --from=builder /opt/app-root/src/bin/generator /
ENTRYPOINT [ "/generator" ]
