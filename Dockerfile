FROM golang:alpine3.13 AS build-env

WORKDIR /src

ADD . /src

RUN cd /src && go build -o main

FROM alpine
WORKDIR /app
COPY --from=build-env /src/main /app/

ENTRYPOINT [ "./main" ]
