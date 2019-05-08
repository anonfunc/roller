FROM golang:1.12 as build
# RUN apt-get update && apt-get install -y upx-ucl
ADD go.mod go.sum /src/roller/
WORKDIR /src/roller/
RUN go mod download

COPY . /src/roller/
RUN CGO_ENABLED=0 GOOS=linux go build -o roller
# RUN upx roller

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /src/roller/roller /bin/roller
WORKDIR /
ENTRYPOINT ["/bin/roller"]
EXPOSE 3000
