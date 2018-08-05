package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/ego008/selenium"
	"github.com/tebeka/selenium/chrome"
)

func FindLineCount(filename string) int {
	file, _ := os.Open(filename)
	var line int = 0
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line++
	}
	if ok := scanner.Err(); ok != nil {
		return 0
	}
	return line
}

func FindLine(filename string, search int) string {
	file, _ := os.Open(filename)
	var line int = 0
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if line == search {
			return scanner.Text()
		}
		line++
	}

	if ok := scanner.Err(); ok != nil {
		return ""
	}
	return ""
}

func randomString(size int, filename string) string {
	rand.Seed(time.Now().Unix())
	randomInt := rand.Intn(size)
	key := FindLine(filename, randomInt)
	return key
}

func Starting(count3 int, urls string, count2 int, useragents string, count1 int, proxy string, timeout int64) {
	for {
		_proxy := randomString(count1, proxy)
		_useragent := randomString(count2, useragents)
		_url := randomString(count3, urls)
		port, err := pickUnusedPort()
		opts := []selenium.ServiceOption{
			selenium.Output(os.Stderr),
		}
		selenium.SetDebug(false)
		service, err := selenium.NewIeDriverService("chromedriver.exe", port, opts...)
		if err != nil {
			return
		}
		defer func() {
			fmt.Println(service.Stop())
		}()
		caps := selenium.Capabilities{"browserName": "chrome"}
		chromeCaps := chrome.Capabilities{
			Path: "",
			Args: []string{
				"--proxy-server=https://" + _proxy,
				"--no-sandbox=true",
				"--user-agent=" + _useragent,
			},
		}
		caps.AddChrome(chromeCaps)
		wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d", port))
		if err != nil {
			return
		}
		defer func() {
			wd.Quit()
		}()

		if err := wd.Get(_url); err != nil {
			return
		}
		fmt.Println(wd.WindowHandles())
	}
	time.Sleep(time.Second * time.Duration(timeout))
}

func pickUnusedPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	port := l.Addr().(*net.TCPAddr).Port
	if err := l.Close(); err != nil {
		return 0, err
	}
	return port, nil
}

func main() {
	// Parse cmdline arguments using flag package
	timeout := flag.Int64("t", 240, "Timeout Duration (second type)")
	proxy := flag.String("p", "proxy.txt", "Proxy file ip:port")
	urls := flag.String("u", "url.txt", "Url File")
	useragents := flag.String("a", "agents.txt", "User Agent File")
	flag.Parse()

	fmt.Println("Agent is starting...")
	fmt.Print("Url file checking...")
	if _, err := os.Stat(*urls); os.IsNotExist(err) {
		fmt.Println(*urls, "not found!")
		os.Exit(1)
	}
	fmt.Print("OK\n")
	fmt.Print("Proxy file checking...")
	if _, err := os.Stat(*proxy); os.IsNotExist(err) {
		fmt.Println(*proxy, "not found!")
		os.Exit(1)
	}
	fmt.Print("OK\n")
	fmt.Print("User agent file checking...")
	if _, err := os.Stat(*useragents); os.IsNotExist(err) {
		fmt.Println(*useragents, "not found!")
		os.Exit(1)
	}
	fmt.Print("OK\n")

	count1 := FindLineCount(*proxy)
	count2 := FindLineCount(*useragents)
	count3 := FindLineCount(*urls)
	if count1 == 0 {
		fmt.Println(*proxy, "data is not found!")
		os.Exit(1)
	}
	if count2 == 0 {
		fmt.Println(*useragents, "data is not found!")
		os.Exit(1)
	}
	if count3 == 0 {
		fmt.Println(*urls, "data is not found!")
		os.Exit(1)
	}

	go Starting(count3, *urls, count2, *useragents, count1, *proxy, *timeout)
	var input string
	fmt.Scanln(&input)
}
