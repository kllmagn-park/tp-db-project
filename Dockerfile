FROM golang:latest AS builder

WORKDIR /app

COPY . ./
RUN GOAMD64=v3 go build -ldflags "-w -s" ./cmd/main.go

FROM ubuntu:20.04

ENV DEBIAN_FRONTEND=noninteractive 
RUN apt-get -y update && apt-get install -y postgresql-12 && rm -rf /var/lib/apt/lists/*
RUN apt-get -y update && apt-get install -y tzdata
RUN ln -snf /usr/share/zoneinfo/Russia/Moscow /etc/localtime && echo Russia/Moscow > /etc/timezone

USER postgres

RUN /etc/init.d/postgresql start && \
  psql --command "CREATE USER root WITH SUPERUSER PASSWORD 'root';" && \
  createdb -O root dbms && \
  /etc/init.d/postgresql stop

EXPOSE 5432

USER root

WORKDIR /cmd

COPY ./db/db.sql ./db.sql

COPY --from=builder /app/main .

EXPOSE 5000
ENV PGPASSWORD root
CMD service postgresql start && psql -h localhost -d dbms -U root -p 5432 -a -q -f ./db.sql && ./main