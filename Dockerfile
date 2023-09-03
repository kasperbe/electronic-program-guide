FROM golang:1.20-alpine AS build

RUN adduser service-user;echo 'service-user'

WORKDIR /app

COPY . ./

RUN go build -o /epg-service

FROM scratch

ENV TZ Europe/Copenhagen

WORKDIR /

COPY --from=build /epg-service /epg-service

USER 1001
EXPOSE 8080

ENTRYPOINT ["/epg-service"]
