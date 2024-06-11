FROM golang:1.22

ENV KEY=jplkey
ENV SECRET_KEY=jplsecretkey

RUN mkdir -p /app

COPY . /app

RUN chmod -x /app/main.go

RUN chmod -x /app/cmd/server.go

WORKDIR /app

RUN go install

CMD ["jpl","server"]