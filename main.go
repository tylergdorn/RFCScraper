package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

func main() {
	viewFlag := flag.Bool("view", false, "Used to print one document to stdout")
	flag.Parse()
	if len(os.Args) == 1 {
		fmt.Println("Usage: RFCScraper [start] [end] where start and end are valid RFCs")
		fmt.Println("Example: RFCScraper 100 120")
		fmt.Println("You can also use --view to view one RFC")
		fmt.Println("Example: RFCScraper --view 1000")
		os.Exit(0)
	}
	if *viewFlag {
		dl, err := strconv.Atoi(flag.Args()[0])
		if err != nil {
			fmt.Println("Bad usage, must be number")
		}
		view(dl)
	} else {
		start := os.Args[1:2][0]
		end := os.Args[2:3][0]
		startInt, err := strconv.Atoi(start)
		endInt, err2 := strconv.Atoi(end)
		// lazy error checking
		if err != nil || err2 != nil {
			fmt.Println("input should be a number corresponding to a RFC.")
		} else {
			// make sure they use it right
			if start > end {
				fmt.Println("Start should be bigger than end. i.e.: RFCScraper 10 12")
			} else {
				downloadRange(startInt, endInt)
			}
		}
	}
}

// downloads a rfc corresponding to the number provided
func download(number int) error {
	content, err := getRFC(number)
	// get the response
	if err != nil {
		return err
	}
	// if it's not an error, make the directory and write out the file
	_ = os.Mkdir("./rfc", 0777)
	file, err := os.Create("./rfc/" + strconv.Itoa(number) + ".txt")
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write([]byte(content))
	return nil
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

// downloads a range of RFCs from start to end.
func downloadRange(start int, end int) {
	i := start
	numDownloads := end - start + 1 // it's inclusive so we want one more
	var waitGroup sync.WaitGroup
	waitGroup.Add(numDownloads)
	// a channel to handle all errors that come through
	errch := make(chan error, numDownloads)
	for i <= end {
		// spawn goroutines for each file
		go func(w *sync.WaitGroup, number int) {
			defer w.Done()
			err := download(number)
			if err != nil {
				errch <- err
			} else {
				errch <- nil
			}
		}(&waitGroup, i)
		i++
	}
	// wait for em
	waitGroup.Wait()
	for i := start; i < numDownloads; i++ {
		err := <-errch
		if err != nil {
			// print errors if any
			fmt.Println(err)
		}
	}
}
