FROM redis:latest

RUN apt-get -y update
RUN apt-get -y upgrade
RUN apt-get -y install curl gcc make

WORKDIR /code
RUN curl -s -L -O "https://github.com/repsheet/redis_module/archive/0.0.2.tar.gz"
RUN tar xzf 0.0.2.tar.gz

WORKDIR /code/redis_module-0.0.2
RUN make
RUN mkdir -p /lib
RUN cp src/repsheet.so /lib/

EXPOSE 6379

CMD ["redis-server", "--loadmodule", "/lib/repsheet.so"]
