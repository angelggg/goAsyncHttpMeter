package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Counter struct {
	Reads int
	Errors int
	Duration time.Duration
}
var iterations = 10
var counter1 Counter
var counter2 Counter


func main () {

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Web 1")
	web1, _ := reader.ReadString('\n')
	web1 = strings.Replace(web1, "\n", "", -1)
	web1 = requestPrefixes(web1)
	fmt.Println("Web 2")
	web2, _ := reader.ReadString('\n')
	web2 = strings.Replace(web2, "\n", "", -1)
	web2 = requestPrefixes(web2)
	fmt.Print(web1 + " " + web2)
	fmt.Println("How long would you like to test?")
	durationInput , _ := reader.ReadString('\n')
	durationInput = strings.Replace(durationInput, "\n", "", -1)
	conv, _ := strconv.Atoi(durationInput)
	dur := time.Second * time.Duration(conv)
	c1 := make(chan time.Duration)
	c2 := make(chan time.Duration)

	go func (){
		for i := 0 ; i < iterations ; i++{
			c1 <- connectTo(web1, 1)
			time.Sleep(time.Second )
		}
		close(c1)
	}()

	go func (){
		for i := 0 ; i < iterations ; i++{
			c2 <- connectTo(web2, 2)
			time.Sleep(time.Second )
		}
		close(c2)
	}()

	loop:
		for timeout := time.After(dur); ; {
		select{
			case m1 := <- c1:
				if m1 != time.Duration(0){
					counter1.Reads += 1
					counter1.Duration += m1
				}	else {
					counter1.Errors += 1
					time.Sleep(time.Second)
				}
				fmt.Println(counter1)

			case m2 := <- c2:
					if m2 != time.Duration(0){
						counter2.Reads += 1
						counter2.Duration += m2
					}	else {
						counter2.Errors += 1
						time.Sleep(time.Second)

					}
					fmt.Println(counter2)

			case <- timeout:
				break loop
		}
	}
	getInfo(web1, web2)
}

func requestPrefixes (url string) string {
	if strings.HasPrefix(url, "http://"){
		return url
	} else if strings.HasPrefix(url, "https://"){
		return url
	} else {
		return "https://" + url
	}
}
func getInfo(web1 string, web2 string){
	fmt.Println("Web1 " + web1 + ":")
	hits1 := strconv.Itoa(counter1.Reads)
	mean_ct1 := counter1.Duration / time.Duration(counter1.Reads)
	errors1 := strconv.Itoa(counter1.Errors)
	fmt.Println("ERRORS " + errors1 + "|||")
	fmt.Println("Mean connection time " + mean_ct1.String() + "|||")
	fmt.Println("Total connections " + hits1 + "|||")
	fmt.Println("--------------------------------------")

	fmt.Println("Web2 " + web2 + ":")
	hits2 := strconv.Itoa(counter2.Reads)
	mean_ct2 := counter2.Duration / time.Duration(counter2.Reads)
	errors2 := strconv.Itoa(counter2.Errors)
	fmt.Println("ERRORS " + errors2 + "|||")
	fmt.Println("Mean connection time " + mean_ct2.String()  + "|||")
	fmt.Println("Total connections " + hits2 + "|||")
}

func connectTo(url string, web int) time.Duration{
	start := time.Now()
	response , err:= http.Get(url)
	if err == nil {

		fmt.Println("-- PING TO " + url + " -- WITH STATUS " + response.Status)
		if response.StatusCode != 200{
			fmt.Println("STATUS RECEIVED: " + response.Status)
			return time.Duration(0)
		}
		return time.Since(start)
	}
	fmt.Println("ERROR " + err.Error())
	return time.Duration(0)
}
