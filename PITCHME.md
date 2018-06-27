## Web Application Threat Intelligence

Aaron Bedra, Chief Scientist, Jemurai  

[@abedra](https://twitter.com/abedra)  

[keybase.io/abedra](https://keybase.io/abedra)

---

## Agenda

@ul

- Setup
- Adaptive Security
- The Repsheet Environment
- Building a Processor
- Wrap-Up

@ulend

---

## Setup

@fa[arrow-down]

+++

### Build and run the environment

```sh
docker compose up --build -d
```

---

## Adaptive Security

@fa[arrow-down]

+++

## Security and Trust are fluid

+++

### Repsheet gives you more fluid options for dealing with bad actors

+++

## How does it work?

+++

## TODO repsheet image

+++

## Components

@ul

- hiredis (Redis C adapter)
- librepsheet (redis cache layer and core logic)
- repsheet-nginx (NGINX module that provides the control plane)

@ulend

+++

## Let's start with configuration

+++?code=repsheet/nginx.conf

@[8-8]
@[9-9]
@[10-10]
@[11-11]
@[13-13]
@[14-14]
@[15-15]
@[16-16]
@[18-18]
@[19-19]
@[20-20]
@[26-26]
@[27-27]

+++

### A full configuration is provided in the `repsheet-nginx` repository

+++

## Let's see what's going on under the hood

+++

### Take a look at the log file

```sh
docker exec -it sample-app tail -f /go/src/app/logs/app.log
172.19.0.3 - - [10/Jun/2017:13:59:14 +0000] "GET / HTTP/1.0" 200 1371
```

+++

### Pretend we are malicious

```sh
docker run -it --network threat_intel                     \
               --link repsheet-redis:repsheet-redis redis \
               redis-cli -h repsheet-redis
repsheet-redis:6379> set 172.19.0.1:repsheet:ip:blacklisted manual
OK
```

+++

### We can see that things have changed

```sh
curl localhost:8888
```

```html
<html>
  <head><title>403 Forbidden</title></head>
  <body bgcolor="white">
    <center><h1>403 Forbidden</h1></center>
    <hr><center>nginx/1.13.1</center>
  </body>
</html>
```

+++

### We can also change our mind

```fundamental
repsheet-redis:6379> set 172.19.0.1:repsheet:ip:whitelisted manual
OK
curl localhost:8888
-- Normal app output --
```

+++

### Let's take a look at the logs

```fundamental
docker exec -it repsheet tail -f /usr/local/nginx/logs/error.log
```

```fundamental
2017/06/10 14:07:22 [error] 6#0: *102 [Repsheet] - IP 172.19.0.1
  was blocked by repsheet. Reason: manual, client: 172.19.0.1,
  server: , request: "GET / HTTP/1.1", host: "localhost:8888"
2017/06/10 14:09:54 [error] 6#0: *103 [Repsheet] - IP 172.19.0.1
  is whitelisted by repsheet. Reason: manual, client: 172.19.0.1,
  server: , request: "GET / HTTP/1.1", host: "localhost:8888"
```

+++

### We can also be unsure

```text
repsheet-redis:6379> flushdb
OK
repsheet-redis:6379> set 127.19.0.1:repsheet:ip:marked manual
OK
curl localhost:8888
```

```html
<div>
  <p>Actor is marked</p>
</div>
```

+++

### But under the hood

```text
2017/06/10 14:14:20 [error] 6#0: *105 [Repsheet] - IP 172.19.0.1
  was found on repsheet. Reason: manual, client: 172.19.0.1, server: ,
  request: "GET / HTTP/1.1", host: "localhost:8888"
```

+++

### The repsheet module passes a header to the upstream application

+++

### The header indicates a marked actor

+++

### Example

```go
func repsheetHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter,
                                 r *http.Request) {
        if r.Header["X-Repsheet"] != nil {
            context.Set(r, "repsheet", true)
        }
        next.ServeHTTP(w, r)
    })
}
```

+++

### The sample application will show a message

+++

### That's a bad idea

+++

## Lab: Make a better marked response

+++

### Running the solution

```sh
git co solutions/lab1
docker-compose down
docker-compose up --build -d
```

+++

### You will see a captcha if you come from a marked host

+++

### Working with the repsheet cache can be easier

+++

### Example

```sh
cd cli
docker build -t cli .
docker run -it --network threat_intel                            \
               --link repsheet-redis:repsheet-redis cli repsheet \
               -host repsheet-redis -list
blacklisted actors
whitelisted actors
marked actors
  172.19.0.1
```

---

## Automation

@fa[arrow-down]

+++

### Everything you've done so far is manual

+++

### What we have is an improvement, but it lacks automation

+++

### Let's build some

+++

### Before we automate, we need data

+++

### Start by simulating a login attack

```sh
cd pester
docker build -t pester .
docker run --network threat_intel \
           -e TARGET=repsheet     \
           -e PORT=80             \
           pester
```

```text
Starting login attack on http://repsheet:80/
The password is: P4$$w0rd!
```