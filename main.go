package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	arg := os.Args[1:2]
	intarg, err := strconv.Atoi(arg[0])
	if err != nil {
		fmt.Println("input should be a four digit number corresponding to a RFC.")
	} else {
		view(intarg)
	}
}

func view(number int) {
	str, err := getRFC(number)
	if err != nil {
		panic(err)
	}
	fmt.Println(str)
}

func getRFC(number int) (string, error) {
	// build up our link
	var stringbuilder strings.Builder
	stringbuilder.WriteString("https://www.ietf.org/rfc/rfc")
	stringbuilder.WriteString(strconv.Itoa(number))
	stringbuilder.WriteString(".txt")

	// get the final string
	link := stringbuilder.String()

	// get
	resp, err := http.Get(link)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
