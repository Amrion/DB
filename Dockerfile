FROM golang:1.13 AS build

ADD . /app
WORKDIR /app
RUN go build ./cmd/forum/main.go

FROM ubuntu:20.04

RUN apt-get -y update && apt-get install -y tzdata

ENV PGVER 12

COPY db/db.sql /

RUN apt-get -y update && apt-get install -y postgresql-$PGVER

USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -E UTF8 -O docker docker &&\
    /etc/init.d/postgresql stop

EXPOSE 5432

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

WORKDIR /usr/src/app

COPY . .
COPY --from=build /app/main .

EXPOSE 5000

ENV PGPASSWORD docker

CMD service postgresql start && psql -h localhost -d docker -U docker -p 5432 -a -q -f ./db/db.sql && ./main
