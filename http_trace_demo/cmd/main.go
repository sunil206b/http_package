package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptrace"
	"time"
)

func main() {
	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	if err != nil {
		log.Fatal(err)
	}

	t0 := time.Now()
	var getConn, dnsStart, dnsDone, gotConn, gotFirstResByte time.Time

	// Create trace info
	trace := &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			getConn = time.Now()
			fmt.Printf("GetConn(%s) %dms \n", hostPort, getConn.Sub(t0).Milliseconds())
		},
		DNSStart: func(info httptrace.DNSStartInfo) {
			dnsStart = time.Now()
			fmt.Printf("DNSStart(%+v) %dms \n", info, dnsStart.Sub(t0).Milliseconds())
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			dnsDone = time.Now()
			fmt.Printf("DNSDone(%+v) %dms \n", info, dnsDone.Sub(t0).Milliseconds())
		},
		ConnectStart: func(network, addr string) {
			fmt.Printf("ConnectStart(%s, %s)\n", network, addr)
		},
		ConnectDone: func(network, addr string, err error) {
			fmt.Printf("ConnectDone(%s, %s, %v)\n", network, addr, err)
		},
		GotConn: func(gci httptrace.GotConnInfo) {
			gotConn = time.Now()
			fmt.Printf("GotConn(%+v) %d ms \n", gci, gotConn.Sub(t0).Milliseconds())
		},
		GotFirstResponseByte: func() {
			gotFirstResByte = time.Now()
			fmt.Printf("GotFirstResponseByte %dms \n", gotFirstResByte.Sub(t0).Milliseconds())
		},
		PutIdleConn: func(err error) {
			fmt.Printf("PutIdleConn(%v)\n", err)
		},
	}

	//Create trace context
	ctx := httptrace.WithClientTrace(req.Context(), trace)

	//Attach trace context to request
	req = req.WithContext(ctx)

	fmt.Println("# Request to example.com")
	fmt.Println("")

	// HTTP Request/Response - Trace Events
	// DNS Request - DNSStart
	// DNS Response - DNSDone
	// TCP Create Connection - ConnectStart/ConnectDone
	// Write to the TCP connection - WroteHeaders/WroteRequest
	// Read from TCP Connection - GotFirstResponseByte
	// Close TCP connection


	t0 = time.Now()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	// read the whole body and close so that the underlying TCP conn is re-used
	io.Copy(ioutil.Discard, res.Body)
	res.Body.Close()

	fmt.Println("")
	fmt.Println("# Another Request to example.com")
	fmt.Println("")
	t0 = time.Now()
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
}