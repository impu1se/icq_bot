FROM golang:1.13

RUN apt-get update && apt-get install mplayer -y

RUN mkdir /movie_to_gif_bot
ADD . /movie_to_gif_bot/
WORKDIR /movie_to_gif_bot

RUN go mod download
RUN go build -o movie_to_gif_bot cmd/movie_to_gif_bot/main.go

CMD ["/movie_to_gif_bot/movie_to_gif_bot"]

