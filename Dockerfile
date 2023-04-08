FROM golang:1.20-alpine
VOLUME /app/stori-challenge/csv
WORKDIR /app/stori-challenge
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o ./out/stori-challenge ./main/main.go
VOLUME /app/stori-challenge/csv
WORKDIR /app/stori-challenge/csv/incoming
WORKDIR /app/stori-challenge/csv/archive
WORKDIR /app/stori-challenge
COPY resources/20230406_GHM54789345.csv /app/stori-challenge/csv/sample/
CMD ["/app/stori-challenge/out/stori-challenge"]