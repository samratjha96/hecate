FROM golang

WORKDIR /app

COPY . .

ENV GOPROXY=direct

ENV DATABASE_URL="postgres://admin:password@db:5432/hecate?sslmode=disable"

RUN go build -o hecate ./main.go

EXPOSE 3000

CMD [ "./hecate" ]
