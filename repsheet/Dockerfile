FROM centos:latest

RUN yum -y install epel-release
RUN yum -y update
RUN yum -y install gcc make pkg-config pcre-devel zlib-devel hiredis hiredis-devel

WORKDIR /code
RUN curl -s -L -O "https://github.com/repsheet/repsheet-nginx/archive/5.0.2.tar.gz"
RUN tar xzf 5.0.2.tar.gz

RUN curl -s -L -O "http://nginx.org/download/nginx-1.15.7.tar.gz"
RUN tar xzf nginx-1.15.7.tar.gz
WORKDIR /code/nginx-1.15.7
RUN ./configure --add-module=../repsheet-nginx-5.0.2
RUN make
RUN make install

COPY nginx.conf /usr/local/nginx/conf

EXPOSE 80

CMD ["/usr/local/nginx/sbin/nginx"]
