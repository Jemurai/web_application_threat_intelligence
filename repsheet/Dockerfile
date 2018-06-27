FROM centos:latest

RUN yum groupinstall "Development Tools" -y
RUN yum install pcre-devel zlib-devel check-devel -y

WORKDIR /code

RUN git clone https://github.com/redis/hiredis
WORKDIR /code/hiredis
RUN make
RUN PREFIX=/usr LIBRARY_PATH=lib64 make install

WORKDIR /code
RUN curl -s -L -O "https://github.com/repsheet/librepsheet/archive/6.2.1.tar.gz"
RUN tar xzf 6.2.1.tar.gz
WORKDIR /code/librepsheet-6.2.1
RUN ./autogen.sh
RUN ./configure --prefix=/usr --libdir=/usr/lib64
RUN make
RUN make install

WORKDIR /code
RUN curl -s -L -O "https://github.com/repsheet/repsheet-nginx/archive/4.3.0.tar.gz"
RUN tar xzf 4.3.0.tar.gz

RUN curl -s -L -O "http://nginx.org/download/nginx-1.13.1.tar.gz"
RUN tar xzf nginx-1.13.1.tar.gz
WORKDIR /code/nginx-1.13.1
RUN ./configure --add-module=../repsheet-nginx-4.3.0
RUN make
RUN make install

COPY nginx.conf /usr/local/nginx/conf

EXPOSE 80

CMD ["/usr/local/nginx/sbin/nginx"]
