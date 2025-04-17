FROM golang:1.24

WORKDIR /urlShortener

COPY . .

ENV CONFIG_PATH=/urlShortener/config/local.yaml

WORKDIR /urlShortener/cmd/urlShortener
RUN GOOS=linux go build -o urlShortener

EXPOSE 8081

CMD ["./urlShortener"]

# docker build -f .\build\Dockerfile -t alaiv/shortener:latest .
# docker run --rm -p 8081:8081 alaiv/shortener:latest