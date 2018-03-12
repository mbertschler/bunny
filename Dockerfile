FROM golang:1.9.3-alpine3.7 as builder

COPY . /go/src/github.com/mbertschler/bunny

RUN go install github.com/mbertschler/bunny

# ----------------------------------------
FROM alpine:3.7  

RUN apk --no-cache add ca-certificates

COPY --from=builder /go/bin/bunny /bunny
COPY ./js /js

ENV BUNNY_ROOT=/

EXPOSE 3080

CMD ["/bunny"]
