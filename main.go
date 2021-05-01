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
//	"time"
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




func makeConnection(url string, redirect bool, c chan Info) {

	//fmt.Printf("%s[*]%s Starting the petiton to %s%s%s", Yellow, Reset,Purple, url,Purple)

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

	status := resp.StatusCode

	//log.Printf("StatusCode: %d\nStatus Text: %s", status, http.StatusText(status))

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	bod := string(body)

	charas := len(bod)

	lines := strings.Count(bod, "\n")
	lines = lines + 1

	finalurl := resp.Request.URL.String()

	c <- Info{status, charas, lines, finalurl}
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

func printHeader () {
	fmt.Println("=============================================================================")
	fmt.Println("ID	STATUS		Charas		Lines		URL")
	fmt.Println("=============================================================================")
}


func PrintInfo(ch chan Info, index, hc int) {
	//fmt.Printf("%s[*]%s PETITION INFORMATION\n  -Status: %d %s\n  -Characters: %d\n  -Lines: %d  -URL: %s",Blue, Reset, status, http.StatusText(status), charas, lines, finalurl)
	var statusColor = Yellow

	c := <- ch

	if c.status == hc {
		return
	}

	if c.status > 199 && c.status < 300 {
		statusColor = Green
	} else if c.status > 299 && c.status < 400{
		statusColor = Blue
	} else if c.status > 399 && c.status < 500{
		statusColor = Red
	}

	var factor = ""

	if c.status > 199 && c.status < 300 {
		factor = "	"
	} else if c.status == 301 {
		fmt.Printf("%d	%s%d Moved%s%s	%dCh		%dL	%s%s%s", index, statusColor, c.status,factor, Reset, c.charas, c.lines, Blue, c.finalurl, Reset)
		fmt.Println()
		return
	}

	fmt.Printf("%d	%s%d %s%s%s	%dCh		%dL	%s%s%s", index, statusColor,
	 c.status, http.StatusText(c.status), factor, Reset, c.charas, c.lines, Blue, c.finalurl,Reset)

	fmt.Println()

}


func main(){

	url := flag.String("url", "https://en.wikipedia.org/wiki/%22Hello,_World!%22_program", "The url you want to do the attack to :)")

	redirect := flag.Bool("redirect", false, "This will tell the fuzzing to follow redirects")

	dict := flag.String("dict", "nodict", "The dictionary with the payloads that you want to use in your attack")

	hc := flag.Int("hc", 0, "This will hide the code that you want for not to show it on the screen")

	flag.Parse()

	c := make(chan Info)

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

	//PrintInfo(status, charas, lines, finalurl)

	words := getWords(*dict)

	printHeader()

	for index, word := range words {
		//fmt.Println("Index: %d	Word: %s",index, word)
		xurl := clean(*url, word)

		go makeConnection(xurl, *redirect, c)

		PrintInfo(c, index, *hc)
	}

	//fmt.Printf("The words are:")
	//fmt.Printf(words[100])
	//fmt.Println()
}

type Info struct {
	status int
	charas int
	lines int
	finalurl string
}
