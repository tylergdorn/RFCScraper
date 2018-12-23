package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

func main() {
	arg := os.Args[1:2]
	intarg, err := strconv.Atoi(arg[0])
	if err != nil {
		fmt.Print(intarg)
		fmt.Println("input should be a four digit number corresponding to a RFC.")
	} else {
		downloadRange(1000, 1200)
	}
}

// downloads a rfc corresponding to the number provided
func download(number int) error {
	content, err := getRFC(number)
	// get the response
	if err != nil {
		return err
	} else {
		// if it's not an error, make the directory and write out the file
		_ = os.Mkdir("./rfc", 0777)
		file, err := os.Create("./rfc/" + strconv.Itoa(number))
		if err != nil {
			return err
		}
		defer file.Close()
		file.Write([]byte(content))
		return nil
	}
}

// prints an rfc out to console
func view(number int) {
	str, err := getRFC(number)
	if err != nil {
		panic(err)
	}
	fmt.Println(str)
}

// returns the content of the given RFC number.
// takes an int that corresponds to an rfc
func getRFC(number int) (string, error) {
	// build up our link
	var stringbuilder strings.Builder
	stringbuilder.WriteString("https://www.ietf.org/rfc/rfc")
	stringbuilder.WriteString(strconv.Itoa(number))
	stringbuilder.WriteString(".txt")

	// get the final string
	link := stringbuilder.String()

	// get from our built link
	resp, err := http.Get(link)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// check if it exists
	if resp.StatusCode != 200 {
		return "", errors.New("No RFC number: " + strconv.Itoa(number) + "!")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func downloadRange(start int, end int) {
	i := start
	numDownloads := end - start
	var waitGroup sync.WaitGroup
	waitGroup.Add(numDownloads)
	errch := make(chan error, numDownloads)
	for i < end {
		go func(w *sync.WaitGroup) {
			// r := rand.Intn(10)
			// time.Sleep(time.Duration(r) * time.Microsecond)
			err := download(i)
			if err != nil {
				errch <- err
			}
			errch <- nil
			defer w.Done()
		}(&waitGroup)
		i++
	}
	waitGroup.Wait()
	for i := start; i < numDownloads; i++ {
		err := <-errch
		if err != nil {
			fmt.Println(err)
		}
	}
}
