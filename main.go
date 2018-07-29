package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	urlpkg "net/url"
	"os"
	"time"

	"github.com/headzoo/surf"
)

func GetCookie(url string, pref string, visitor string, ysc string) *cookiejar.Jar {
	jar, _ := cookiejar.New(nil)
	var cookies []*http.Cookie
	firstCookie := &http.Cookie{
		Name:   "PREF",
		Value:  pref,
		Path:   "/",
		Domain: ".youtube.com",
	}

	cookies = append(cookies, firstCookie)

	secondCookie := &http.Cookie{
		Name:   "VISITOR_INFO1_LIVE",
		Value:  visitor,
		Path:   "/",
		Domain: ".youtube.com",
	}

	cookies = append(cookies, secondCookie)

	thirdCookie := &http.Cookie{
		Name:   "YSC",
		Value:  ysc,
		Path:   "/",
		Domain: ".youtube.com",
	}

	cookies = append(cookies, thirdCookie)

	fourthCookie := &http.Cookie{
		Name:   "GPS",
		Value:  "1",
		Path:   "/",
		Domain: ".youtube.com",
	}

	cookies = append(cookies, fourthCookie)

	cookieURL, _ := urlpkg.Parse(url)

	jar.SetCookies(cookieURL, cookies)
	return jar
}

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

func Start(count3 int, urls string, count2 int, useragents string, count1 int, proxy string, timeout int64) {
	for {
		_proxy := randomString(count1, proxy)
		_useragent := randomString(count2, useragents)
		_url := randomString(count3, urls)
		bow := surf.NewBrowser()
		bow.SetUserAgent(_useragent)
		jar := GetCookie(_url, "f1=50000000", "Wb4lvmoVRNI", "tm6XNnZbjHI")
		bow.SetCookieJar(jar)
		px, _ := urlpkg.Parse(_proxy)
		transport := http.Transport{
			Proxy:           http.ProxyURL(px),
			TLSClientConfig: &tls.Config{},
		}
		bow.SetTransport(&transport)
		err := bow.Open(_url)
		time.Sleep(time.Duration(time.Second * time.Duration(timeout)))
		if err != nil {
			fmt.Println(os.Stderr, "open url error", err)
		}
		fmt.Println("\n Useragent:", _useragent, "Proxy:", _proxy, "Url:", _url, "Response:", bow.ResponseHeaders(), "\n")
	}
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

	go Start(count3, *urls, count2, *useragents, count1, *proxy, *timeout)
	var input string
	fmt.Scanln(&input)
}
