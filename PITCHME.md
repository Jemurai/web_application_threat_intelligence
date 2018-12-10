## Web Application Threat Intelligence

Aaron Bedra, Chief Scientist, Jemurai  

[@abedra](https://twitter.com/abedra)  

[keybase.io/abedra](https://keybase.io/abedra)

---

## Agenda

@ul

- Setup
- Adaptive Security
- Identifying Threats
- Indicators of Attack vs Indicators of Compromise
- Containment of Malicious Actors
- Device Fingerprinting
- Wrap-Up

@ulend

---

## Setup

@fa[arrow-down]

+++

### Build and run the environment

```sh
docker-compose up --build -d
```

---

## Adaptive Security

---

## Identifying Threats

@fa[arrow-down]

+++

### How do you know a request or an actor is malicious?

+++

### That's a loaded question

+++

### Let's start by generating some data

+++

### Simulated login attack

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

+++

### Grab the logs from our sample application

```sh
cd ../processor
docker cp sample-app:/go/src/app/logs/app.log .
```

+++

### What do you see?

+++

### It looks pretty easy to disect

+++

### So let's build something to do that

+++?code=processor/main.go&lang=golang

@[42-54]

+++

### Give it a try

```sh
docker build -t processor .
docker run --network threat_intel processor
```

```sh
Blacklisting 172.19.0.5. Threshold: 10, Actual: 31
```

+++

### What does this solution provide?

+++

### A flexible way to ask questions about request traffic

+++

### You can build any number of detectors on top of this model

+++

### There are commercial tools available that help identify basic threats like this

+++

### But for the more advanced detectors it's better to have your own processors

+++

### If we have time later, we will build a more advanced detector

+++

### Wrap-Up

---

## Indicators of Attack vs Indicators of Compromise

@fa[arrow-down]

+++

### So far we have only asked questions about attacks

+++

### These are referred to as Indicators of Attack

+++

### What do you do if you know you are being attacked?

@ul

- Let it happen
- Block the actor
- Report the actor
- Confuse the actor

@ulend

+++

### What if an attack is successful?

+++

### What do you do then?

+++

### Take a look back at your application logs

```nohighlight
2018/12/10 01:30:10 Login failed for admin@example.com
2018/12/10 01:30:10 Login failed for admin@example.com
2018/12/10 01:30:10 Login failed for admin@example.com
2018/12/10 01:30:10 Login failed for admin@example.com
2018/12/10 01:30:10 Successfull login for admin@example.com
```

+++

### Notice anything about the last entry?

+++

### The 302 is a successful login

+++

### If an actor trips the brute force detector then logs in, what does that mean?

+++

### This is an Indicator of Compromise

+++

### Lab: Detect brute force IoC

+++

### Now that we can detect an IoC, what do we do?

+++

### In some cases you can automate incident response

+++

### In the case of an account takeover, you can disable the account

+++

### And force the reset of the users password

+++

### If there are payment instruments, delete them

+++

### Responding to an IoC is about preventing loss

+++

### Loss doesn't occur until an asset is misused

+++

### Wrap-Up

---

## Containment of Malicous Actors

@fa[arrow-down]

+++

### Now that we are able to identify bad actors, what do we do with them?

+++

### First, you want to stop them from any continued activity

+++

### Depending on how you deploy, this can be difficult

+++

### We're going to assume deployment behind a reverse proxy server

+++

### In this case, NGINX

+++

### Our environment consists of an application, NGINX, and Redis

+++?code=docker-compose.yml&lang=yaml

@[3-9]
@[10-14]
@[15-21]

+++

### Our control environment is already wired up

+++

### But before we start exploring let's understand how it works

+++

### Repsheet gives you more fluid options for dealing with bad actors

+++

## How does it work?

+++

## Components

@ul

- redis_module (custom cache handling)
- repsheet-nginx (NGINX module)

@ulend

+++

## Let's start with configuration

+++?code=repsheet/nginx.conf

@[8-8]
@[9-9]
@[10-10]
@[11-11]
@[12-12]
@[13-13]
@[14-14]
@[15-15]
@[17-17]
@[19-21]
@[23-29]

+++

### A full configuration is provided in the `repsheet-nginx` repository

+++

## Let's see what's going on under the hood

+++

### Take a look at the log file

```sh
docker exec -it sample-app tail -f /go/src/app/logs/app.log
172.18.0.1 - - [10/Dec/2018:13:04:46 +0000] "GET / HTTP/1.0" 200 1545
```

+++

### Pretend we are malicious

```sh
docker exec -it repsheet-redis redis-cli
repsheet-redis:6379> REPSHEET.BLACKLIST 172.18.0.1 manual
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
    <hr><center>nginx/1.15.7</center>
  </body>
</html>
```

+++

### We can also change our mind

```sh
repsheet-redis:6379> REPSHEET.WHITELIST 172.18.0.1 manual
OK
curl localhost:8888
-- Normal app output --
```

+++

### Let's take a look at the logs

```sh
docker exec -it repsheet tail -f /usr/local/nginx/logs/error.log
```

```log
2018/12/10 13:05:55 [error] 6#0: *4 [Repsheet] - IP 172.18.0.1
  was blocked by repsheet. Reason: manual, client: 172.18.0.1,
  server: , request: "GET / HTTP/1.1", host: "localhost:8888"
2018/12/10 13:07:32 [error] 6#0: *5 [Repsheet] - IP 172.18.0.1
  is whitelisted by repsheet. Reason: manual, client: 172.18.0.1,
  server: , request: "GET / HTTP/1.1", host: "localhost:8888"
```

+++

### We can also be unsure

```sh
repsheet-redis:6379> flushdb
OK
repsheet-redis:6379> REPSHEET.MARK 172.18.0.1 manual
OK
curl localhost:8888
```

```html
<div class="g-recaptcha"
  data-sitekey="6LcuXBAUAAAAAJ0vN6S5bzN9nq_Pn5uIxrxEqsBz">
</div>
```

+++

### But under the hood

```log
2018/12/10 13:10:54 [error] 6#0: *7 [Repsheet] - IP 172.18.0.1
  was found on repsheet. Reason: manual, client: 172.18.0.1, server: ,
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

### The sample application will show a captcha

+++

### If you haven't, run the processor to see the complete idea in action

+++

### Wrap-Up

---

## Device Fingerprinting

@fa[arrow-down]

+++

### Fingerprinting helps you keep track of clients

+++

### It provides a ton of useful information you otherwise can't get

+++

### It relies on client side code execution

+++

### This used to be done via flash

+++

### But we all know that doesn't win friends anymore

+++

### So now it's done via JavaScript

+++

### In particular, via a few specific functions

+++

[https://developer.mozilla.org/en-US/docs/Web/API/Navigator](https://developer.mozilla.org/en-US/docs/Web/API/Navigator)

+++

### The navigator object exposes all of the basic information about the client browser

+++

### You can use it to detect if the actor is inconsistent

+++

[https://developer.mozilla.org/en-US/docs/Web/API/WebGL_API/Tutorial](https://developer.mozilla.org/en-US/docs/Web/API/WebGL_API/Tutorial)

+++

### WebGL allows you to execute more advanced code and better determine the browser's fingerprint

+++

### This topic goes quite deep, so be sure to read more

+++

[https://github.com/Valve/fingerprintjs2](https://github.com/Valve/fingerprintjs2)

+++

### Let's give it a try

```js
new Fingerprint2().get(function(result, components) {
  console.log(result) // a hash, representing your device fingerprint
  console.log(components) // an array of FP components
})
```

+++

## Lab: Record device fingerprints

+++

### What did you do with them?

+++

### What kinds of things might you look for given the new information?

+++

### What would you do in the absence of the information?

+++

[https://www.youtube.com/watch?v=SdkKKmL-B_U](https://www.youtube.com/watch?v=SdkKKmL-B_U)

+++

### Wrap-Up

---

## Applying Machine Learning to Threat Intelligence

@fa[arrow-down]

+++

### It's pretty easy to imagine the possible applications of ML to this problem space

+++

### But often times it's a poor choice

+++

### The reality is that both humans and bots are fairly predictable

+++

### And the best solution is to extract functions

+++

### ML can be slow and costly

+++

### So before you reach for it make sure you can't do the same thing with functions and counters

+++

### The best application of ML is in pattern discovery

+++

### Example: Clustering

+++

### Let's look at training a model

```sh
cd fwaf
docker build -t fwaf .
docker run fwaf
```

+++

### Try a scanning tool

```sh
docker run -it frapsoft/nikto \
--link repsheet:repsheet      \
--network threat_intl         \
-host http://repsheet
```

+++

### Optional Lab: Run logs through fwaf to identify malicious requests

+++

### Wrap-Up

---

## Parting thoughts and questions