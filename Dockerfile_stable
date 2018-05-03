# Dockerfile using stable Maxfuzz image
# TODO: FROM Image on docker hub

WORKDIR /root/fuzzer-files
COPY ./fuzzers/ ./
COPY ./.git .
RUN chmod -R 755 /root/fuzzer-files
RUN echo "export GIT_SHA=$(git rev-parse HEAD)" >> /root/fuzzer-files/base/environment
