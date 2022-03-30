FROM golang:1.17

ARG WORKDIR=/src

RUN apt-get update && apt-get install -y \
  git \
  && go install -v \
      golang.org/x/tools/gopls@latest \
  && rm -rf /var/lib/apt/lists/*

ENV LANG=C.UTF-8 TZ=Asia/Tokyo

RUN mkdir -p $WORKDIR

WORKDIR $WORKDIR

VOLUME ["$WORKDIR"]

CMD ["/bin/bash"]
