package main

import (
	"net/http"
	"flag"
	"fmt"
	"log"
	"io/ioutil"
	"strings"
	"os"
)


var Reset  = "\033[0m"
var Red    = "\033[31m"
var Green  = "\033[32m"
var Yellow = "\033[33m"
var Blue   = "\033[34m"
var Purple = "\033[35m"
var Cyan   = "\033[36m"
var Gray   = "\033[37m"
var White  = "\033[97m"




func makeConnection(url string, redirect bool) (status, charas, lines int, finalurl string){

	fmt.Printf("%s[*]%s Starting the petiton to %s%s%s", Yellow, Reset,Purple, url,Purple)

	client := &http.Client{}

	if !redirect {
		client = &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}}
	}

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

	bod := string(body)

	charas = len(bod)

	lines = strings.Count(bod, "\n")
	lines = lines + 1

	finalurl = resp.Request.URL.String()

	return status, charas, lines, finalurl
}


func clean(url, dict string) (xurl, xdict string){

	xurl = strings.ReplaceAll(url, "GO", "")

	return xurl, xdict
}


func main(){

	url := flag.String("url", "https://en.wikipedia.org/wiki/%22Hello,_World!%22_program", "The url you want to do the attack to :)")

	redirect := flag.Bool("redirect", false, "This will tell the fuzzing to follow redirects")

	dict := flag.String("dict", "nodict", "The dictionary with the payloads that you want to use in your attack")

	flag.Parse()

	if *dict == "nodict"{
		fmt.Printf("%s[!]%s There are no payloads for attacking", Red, Reset)
		fmt.Println()
		fmt.Printf("%s[!] %s", Red, Reset)
		os.Exit(1)
	}

	if !strings.Contains(*url, "GO") {
		fmt.Printf("%s[!]%s You must specify where to substitude with the word GO", Red,Reset)
		fmt.Println()
		fmt.Printf("%s[!] %s", Red, Reset)
		os.Exit(1)
	}

	xurl, payloads := clean(*dict, *url)

	status, charas, lines, finalurl := makeConnection(xurl, *redirect)

	fmt.Println()

	fmt.Printf("%s[*]%s PETITION INFORMATION\n  -Status: %d %s\n  -Characters: %d\n  -Lines: %d  -URL: %s",Yellow, Reset, status, http.StatusText(status), charas, lines, finalurl)

	fmt.Println()
}

