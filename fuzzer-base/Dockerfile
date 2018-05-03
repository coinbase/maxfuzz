FROM ubuntu:trusty@sha256:ed49036f63459d6e5ed6c0f238f5e94c3a0c70d24727c793c48fded60f70aa96

MAINTAINER Everest Munro-Zeisberger

WORKDIR /root

########################
# SETUP ENV & VERSIONS #
########################

# Versions:
ENV AFL_VERSION 2.52b
ENV RUBY_VERSION 2.3.3
ENV GO_DEP_VERSION 0.4.1

# Environment Variables:
ENV GOPATH=/root/go/
ENV GOBIN=/root/go/bin

################
# INSTALL DEPS #
################

RUN apt-get update
RUN apt-get install -y software-properties-common
RUN apt-add-repository -y ppa:rael-gc/rvm
RUN apt-add-repository -y ppa:gophers/archive
RUN apt-get update
RUN apt-get install -y \
  git\
  wget\
  gcc\
  autoconf\
  make\
  bison\
  libssl-dev\
  libreadline-dev\
  zlib1g-dev\
  pkg-config\
  gcc\
  clang\
  llvm\
  rvm\
  watch\
  nodejs\
  npm\
  supervisor\
  cmake\
  gdb\
  python-virtualenv\
  golang-1.9-go\
  cython\
  build-essential\
  libgtk2.0-dev\
  libtbb-dev\
  python-dev\
  python-numpy\
  python-scipy\
  libjasper-dev\
  libjpeg-dev\
  libpng-dev\
  libtiff-dev\
  libavcodec-dev\
  libavutil-dev\
  libavformat-dev\
  libswscale-dev\
  libdc1394-22-dev\
  libv4l-dev

###########################
# AFL Compilation & Setup #
###########################

# Download AFL and uncompress
RUN wget http://lcamtuf.coredump.cx/afl/releases/afl-$AFL_VERSION.tgz
RUN tar -xvf afl-$AFL_VERSION.tgz
RUN mv afl-$AFL_VERSION afl

# Inject our own AFL config header file
RUN rm /root/afl/config.h
COPY ./config/afl_config/config.h /root/afl/config.h

# Compile both standard gcc, clang etc as well as afl-clang-fast, used for
# faster & persistent test harnesses. Also build the afl-fuzz binary
RUN cd ~/afl && make
RUN cd ~/afl/llvm_mode && make

# Environment Setup
ENV AFL_I_DONT_CARE_ABOUT_MISSING_CRASHES="1"

# Install py-afl-fuzz (for fuzzing python libraries)
RUN git clone https://github.com/jwilk/python-afl.git
RUN cd python-afl && python setup.py install

# Compile ruby from sources with afl, and setup cflags to access instrumented
# ruby headers (useful for ruby library fuzzing with C harneses)
RUN CC=~/afl/afl-clang-fast /usr/share/rvm/bin/rvm install --disable-binary $RUBY_VERSION
ENV LD_LIBRARY_PATH="LD_LIBRARY_PATH=/usr/share/rvm/rubies/ruby-$RUBY_VERSION/lib"
ENV PATH="/usr/share/rvm/rubies/ruby-$RUBY_VERSION/bin:$PATH"
COPY ./config/ ./config/
RUN PKG_CONFIG_PATH=/usr/share/rvm/rubies/ruby-$RUBY_VERSION/lib/pkgconfig pkg-config --cflags --libs ruby-2.3 > ~/config/afl-ruby-flags

# File structure setup
RUN mkdir ~/fuzz_out
RUN mkdir ~/fuzz_in

############################
#SIDECAR & MONITORING SETUP#
############################

# Setup logging and scripts & install Go libs
RUN mkdir /root/logs
WORKDIR /root/go/src/maxfuzz/fuzzer-base
RUN wget https://github.com/golang/dep/releases/download/v$GO_DEP_VERSION/dep-linux-amd64
RUN mv dep-linux-amd64 /usr/local/bin/dep
ENV PATH=/usr/lib/go-1.9/bin/:$PATH
RUN chmod +x /usr/local/bin/dep
RUN go get -u github.com/dvyukov/go-fuzz/...

# Copy Go files into container & compile binaries
RUN mkdir -p /root/go/src/maxfuzz/fuzzer-base
WORKDIR /root/go/src/maxfuzz/fuzzer-base
COPY ./cmd /root/go/src/maxfuzz/fuzzer-base/cmd
COPY ./internal /root/go/src/maxfuzz/fuzzer-base/internal
COPY ./Gopkg.lock /root/go/src/maxfuzz/fuzzer-base/Gopkg.lock
COPY ./Gopkg.toml /root/go/src/maxfuzz/fuzzer-base/Gopkg.toml
COPY ./Makefile /root/go/src/maxfuzz/fuzzer-base/Makefile
RUN make && make install
WORKDIR /root

# Setup sidecar webapp
COPY ./sidecar ./sidecar
RUN curl -sL https://deb.nodesource.com/setup_8.x | sudo -E bash -
RUN apt-get install -y nodejs
RUN cd sidecar && npm install log

#################################
# FINAL SETUP & COPYING SCRIPTS #
#################################

RUN mkdir fuzzer-files
COPY ./scripts ./scripts
COPY ./fuzzer-files/base ./fuzzer-files/base
RUN chmod -R 755 /root/fuzzer-files
RUN chmod 755 /root/scripts/reproduce_stdin
RUN echo "export GIT_SHA=test_environment" >> /root/fuzzer-files/base/environment
