package misc

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func WelcomeMessage() {
	var (
		Green = "\033[32m"
		Reset = "\033[0m"
	)
	b, err := os.ReadFile("../../assets/ascii.txt")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Welcome to...%s%s%s", Green, b, Reset)
}

func FormatURL(rawURL, rawPath string, values *url.Values) string {
	serverURL, err := url.Parse(rawURL)
	if err != nil {
		log.Fatal(err)
	}

	joined, err := url.JoinPath(serverURL.Path, rawPath)
	if err != nil {
		log.Fatal(err)
	}

	serverURL.Path = joined
	serverURL.RawQuery = values.Encode()

	return serverURL.String()
}

func SplitToString(a []int, sep string) string {
	if len(a) == 0 {
		return ""
	}

	b := make([]string, len(a))
	for i, v := range a {
		b[i] = strconv.Itoa(v)
	}
	return strings.Join(b, sep)
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}
