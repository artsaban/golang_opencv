FROM gocv/opencv

ENV GOPATH /go
COPY . /go/src/gocv/
WORKDIR /go/src/gocv

RUN go get gocv.io/x/gocv
RUN go get github.com/lucasb-eyer/go-colorful
RUN go build ./cmd/main.go
