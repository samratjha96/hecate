FROM golang

RUN mkdir /app

ADD . /app

WORKDIR /app

ENV GOPROXY=direct

RUN go build -o hecate ./main.go

EXPOSE 8080
CMD [ "/app/hecate" ]
