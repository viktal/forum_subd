FROM golang:1.15.2-buster AS build

# Собираем генераторы
WORKDIR /opt/build
RUN go get -v github.com/rubenv/sql-migrate/... && mv $GOHOME/go/bin/sql-migrate .

# Копируем исходный код в Docker-контейнер
COPY . /opt/build

RUN go build ./main.go

FROM ubuntu:20.04 AS release

MAINTAINER Talmaza Viktoria

# Make the "en_US.UTF-8" locale so postgres will be utf-8 enabled by default
RUN apt-get -y update && apt-get install -y tzdata

ENV TZ=Russia/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

#
# Install postgresql
#
RUN apt-get update -y && apt-get install -y postgresql postgresql-contrib

# Run the rest of the commands as the ``postgres`` user created by the ``postgres-$PGVER`` package when it was ``apt-get installed``
USER postgres

# Create a PostgreSQL role named ``docker`` with ``docker`` as the password and
# then create a database `docker` owned by the ``docker`` role.
RUN /etc/init.d/postgresql start &&\
    psql --command "ALTER USER postgres WITH PASSWORD 'docker';" &&\
    createdb -O postgres forum_subd &&\
    /etc/init.d/postgresql stop

# Expose the PostgreSQL port
EXPOSE 5432

# Add VOLUMEs to allow backup of config, logs and databases
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Back to the root user
USER root

# Объявлем порт сервера
EXPOSE 5000

# Собранный ранее сервер
WORKDIR /usr/src/app

COPY ./migrations migrations
COPY ./configs configs
COPY ./dbconfig.yml dbconfig.yml

COPY --from=build /opt/build/main .
COPY --from=build /opt/build/sql-migrate .

CMD service postgresql start && ./sql-migrate up && ./main
