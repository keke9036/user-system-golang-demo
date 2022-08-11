package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"sync"
	"sync/atomic"
	"time"
)

var (
	n      int
	c      int
	random bool
)

func init() {
	flag.IntVar(&n, "n", 0, "Number of requests")
	flag.IntVar(&c, "c", 0, "Number of concurrency")
	flag.BoolVar(&random, "random", false, "random test")
}

func login(c int) []*http.Client {
	var clients []*http.Client
	for i := 0; i < c; i++ {
		cookieJar, err := cookiejar.New(nil)
		if err != nil {
			panic(err)
		}
		client := &http.Client{
			Jar:     cookieJar,
			Timeout: 5 * time.Second,
		}
		password := "rootpwd"
		name := fmt.Sprintf("testu_%v", rand.Intn(10000000))
		json := "{\"username\":\"" + name + "\",\"password\":\"" + password + "\"}"
		fmt.Println("Start login: ", name)

		req, err := http.NewRequest("POST",
			"http://localhost:8080/api/v1/user/login",
			bytes.NewBuffer([]byte(json)))
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
		}
		if resp.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			log.Println("Login failed:" + string(body))
		}

		//bodyStr, _ := ioutil.ReadAll(resp.Body)
		//fmt.Println("body ", bodyStr)
		_, err = io.Copy(ioutil.Discard, resp.Body)
		if err != nil {
			log.Println(err)
			return nil
		}
		resp.Body.Close()
		clients = append(clients, client)
	}
	fmt.Println("All clients login successfully")
	return clients
}

func logout(clients []*http.Client) {
	bodyBuf := &bytes.Buffer{}
	for _, client := range clients {
		objReq, err := http.NewRequest("POST", "http://localhost:8080/api/v1/user/logout", bodyBuf)
		resp, err := client.Do(objReq)
		if err != nil {
			log.Print(err)
		}
		if resp.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			log.Print("Logout failed:" + string(body))
		}
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
	fmt.Println("All clients logout successfully")
}

func main() {
	flag.Parse()
	total := 1
	for i := 0; i < total; i++ {
		Once()
	}
}

func Once() {
	wg := new(sync.WaitGroup)
	wg.Add(c)
	remaining := int32(n)
	var (
		successCnt int32
	)

	// Random users test
	benchTest := func(client *http.Client) {
		fmt.Println("Start benchmark")
		bodyBuf := &bytes.Buffer{}
		for atomic.AddInt32(&remaining, -1) >= 0 {
			objReq, err := http.NewRequest("GET", "http://localhost:8080/api/v1/user/info", bodyBuf)
			resp, err := client.Do(objReq)
			if err != nil {
				log.Print(err)
				continue
			}
			//defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				body, _ := ioutil.ReadAll(resp.Body)
				log.Print("Request failed:", string(body))
				resp.Body.Close()
				continue
			}
			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
			atomic.AddInt32(&successCnt, 1)
		}
		wg.Done()
	}

	var clients []*http.Client
	var start time.Time

	if random {
		fmt.Println("Random user test")
		// Login at first.
		clients = login(c)
		// Logout at the end.
		defer logout(clients)
		start = time.Now()
		for _, client := range clients {
			go benchTest(client)
		}
	} else {
		fmt.Println("Single user test")
		clients = login(1)
		defer logout(clients)
		start = time.Now()
		for i := 0; i < c; i++ {
			go benchTest(clients[0])
		}
	}

	wg.Wait()
	elapsed := time.Since(start)

	fmt.Println("Success request num:", successCnt)
	fmt.Println("Failed request num:", n-int(successCnt))
	fmt.Printf("\tTotal Requests(%v) - Concurrency(%v) - Cost(%s) - QPS(%v/sec)\n",
		n, c, elapsed, math.Ceil(float64(n)/(float64(elapsed)/1000000000)))
}
