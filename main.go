package main

import (
	"net/http"
	"flag"
	"fmt"
	"log"
	"io/ioutil"
	"strings"
	"os"
	"bufio"
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


func getWords(dict string) (words []string) {
	f, err := os.Open(dict)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	s := bufio.NewScanner(f)

	words = make([]string, 0)

	for s.Scan(){
		fmt.Printf("One word is %s", s.Text())
		fmt.Println()

		words = append(words, s.Text())

	}

	err = s.Err()

	if err != nil {
		log.Fatal(err)
	}

	return words

}

func clean(url, word string) (xurl string){

	xurl = strings.ReplaceAll(url, "GO", word)

	return xurl
}

func PrintInfo(status, charas, lines int, finalurl string) {
	//fmt.Printf("%s[*]%s PETITION INFORMATION\n  -Status: %d %s\n  -Characters: %d\n  -Lines: %d  -URL: %s",Blue, Reset, status, http.StatusText(status), charas, lines, finalurl)
	var statusColor = Yellow

	if status > 199 && status < 300 {
		statusColor = Green
	} else if status > 299 && status < 400{
		statusColor = Blue
	} else if status > 399 && status < 500{
		statusColor = Red
	}

	fmt.Printf("       %s%d%s%s        %dW         %dL      %s%s%s", statusColor,
	 status, http.StatusText(status), Reset, charas, lines, Blue, finalurl,Reset)


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

	//xurl := clean(*dict, *url)

	//fmt.Printf("%s[*] %sThe payloads are %s",Yellow, Reset, payloads)

	//fmt.Println()

	//fmt.Printf("%s[*] %sThe url is %s", Yellow, Reset, xurl)

	fmt.Println()

	//status, charas, lines, finalurl := makeConnection(xurl, *redirect)

	fmt.Println()
	
	fmt.Printf("=======================================================================\n")

	fmt.Printf("       STATUS       Characters     Lines     URL\n")

	fmt.Printf("=======================================================================\n")

	//PrintInfo(status, charas, lines, finalurl)

	fmt.Println()

	words := getWords(*dict)


	for index, word := range words {
		//fmt.Println("Index: %d	Word: %s",index, word)

		fmt.Printf("%d", index)

		fmt.Println()

		xurl := clean(*url, word)

		status, charas, lines, finalurl := makeConnection(xurl, *redirect)

		PrintInfo(status, charas, lines, finalurl)
	}

	//fmt.Printf("The words are:")
	//fmt.Printf(words[100])
	//fmt.Println()
}

