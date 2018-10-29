FROM golang:1.11
WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
# RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
# RUN dep ensure -v
RUN go install -v ./...

CMD ["app", "--port" ,":4000"]

