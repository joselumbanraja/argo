FROM golang:1.23-alpine AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /out/webapp .

FROM alpine:3.20
COPY --from=build /out/webapp /webapp
USER 65534
EXPOSE 8080
ENTRYPOINT ["/webapp"]