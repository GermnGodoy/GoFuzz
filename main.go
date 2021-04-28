package main

import (
	"net/http"
	"flag"
	"fmt"
	"log"
	"io/ioutil"
	"strings"
)

func makeConnection(url string) (status, charas, lines int){

	fmt.Printf("[*] Starting the petiton to %s", url)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	resp, _ := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	status = resp.StatusCode

	//log.Printf("StatusCode: %d\nStatus Text: %s", status, http.StatusText(status))

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	charas = len(body)

	bod := string(body)

	lines = strings.Count(bod, "\n")
	lines = lines + 1

	return status, charas, lines
}


func main(){

	url := flag.String("url", "https://en.wikipedia.org/wiki/%22Hello,_World!%22_program", "The url you want to do the attack to :)")

	flag.Parse()

	status, charas, lines := makeConnection(*url)

	fmt.Printf("[*] PETITION INFORMATION\n  -Status: %d %s\n  -Characters: %d\n  -Lines: %d", status, http.StatusText(status), charas, lines)
}

