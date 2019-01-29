FROM golang:1.10

WORKDIR /go/src/github.com/paulloz/bip-boup

COPY . /go/src/github.com/paulloz/bip-boup

RUN go get github.com/bwmarrin/discordgo
RUN go build

CMD ["bip-boup"]
