FROM golang:1.23

LABEL  maintainer "rspepe <rspepeg@gmail.com>"

# air をインストール（ホットリロード用）
RUN go install github.com/cosmtrek/air@latest

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

WORKDIR /src
COPY src .

CMD ["/entrypoint.sh"]
