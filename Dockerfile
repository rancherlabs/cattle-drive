FROM golang:1.22.0-alpine3.19 AS build
RUN apk add -U --no-cache make git bash
COPY . /src/cattle-drive
WORKDIR /src/cattle-drive
RUN ls -l
RUN make build

FROM alpine AS package
COPY --from=build /src/cattle-drive/bin /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/cattle-drive"]
