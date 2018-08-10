FROM golang:rc-alpine

COPY app.go /gobuild/app.go
COPY go.mod /gobuild/go.mod
COPY go.sum /gobuild/go.sum
COPY rds-combined-ca-bundle.pem /gobuild/rds-combined-ca-bundle.pem

RUN apk add --no-cache git

WORKDIR /gobuild
RUN CGO_ENABLED=0 GOOS=linux go build -o /app app.go
#RUN CGO_ENABLED=0 GOOS=linux go build -o /app /gobuild/app.go

FROM scratch
COPY --from=0 /app /app
COPY --from=0 /gobuild/rds-combined-ca-bundle.pem /rds-combined-ca-bundle.pem
EXPOSE 80
CMD ["/app"]

