FROM golang:1.10

WORKDIR /go/src/github.com/paulloz/bip-boup

COPY . /go/src/github.com/paulloz/bip-boup

RUN go get github.com/bwmarrin/discordgo
RUN go get github.com/gojp/kana
RUN go get github.com/ikawaha/kagome/tokenizer
RUN go get github.com/moovweb/gokogiri
RUN go get golang.org/x/text/transform
RUN go get golang.org/x/text/unicode/norm
RUN go build

ENV DISCORD_TOKEN ""

CMD ["bip-boup"]
