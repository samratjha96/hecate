FROM golang

# Install SQLite dependencies
RUN apt-get update && apt-get install -y sqlite3 build-essential

WORKDIR /app

# Create data directory for SQLite
RUN mkdir -p /app/data && chmod 777 /app/data

COPY . .

ENV GOPROXY=direct
ENV SERVER_PORT=8000

# Build with SQLite support
RUN CGO_ENABLED=1 go build -o hecate *.go

EXPOSE 8000

CMD [ "./hecate" ]
