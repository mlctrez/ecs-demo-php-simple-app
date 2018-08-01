FROM golang:1.11beta2-alpine3.8

COPY app.go /gobuild/app.go
COPY go.mod /gobuild/go.mod
COPY go.sum /gobuild/go.sum

RUN apk add --no-cache git

RUN cd /gobuild && CGO_ENABLED=0 GOOS=linux go build -o /app app.go
#RUN CGO_ENABLED=0 GOOS=linux go build -o /app /gobuild/app.go

FROM scratch
COPY --from=0 /app /app
EXPOSE 80
CMD ["/app"]

