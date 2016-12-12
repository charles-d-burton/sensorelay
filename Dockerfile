FROM golang

ADD ./* /go/src/sensorelay/

RUN go get github.com/hashicorp/mdns 

RUN go install sensorelay/

ENTRYPOINT /go/bin/sensorelay

EXPOSE 9899
