FROM golang:1.6

RUN apt-get update
RUN apt-get -y install libvips-dev
RUN apt-get -y install libgsf-1-dev
RUN ldconfig

RUN mkdir -p /go/src/hungry-hippo
WORKDIR /go/src/hungry-hippo

COPY . /go/src/hungry-hippo

RUN go install
RUN chmod 755 wait-for-it.sh

ENV GOOGLE_APPLICATION_CREDENTIALS /go/src/hungry-hippo/treasure-dev.json

CMD ["hungry-hippo", "-queues=fetch_queue", "-use-number", "-uri=redis://redis:6379/"]
