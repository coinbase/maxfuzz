# Dockerfile using most recently built local Maxfuzz image
FROM maxfuzz:latest

WORKDIR /root/fuzzer-files
COPY ./fuzzers/ ./
COPY ./.git .
RUN chmod -R 755 /root/fuzzer-files
RUN echo "export GIT_SHA=$(git rev-parse HEAD)" >> /root/fuzzer-files/base/environment
