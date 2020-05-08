package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var (
	d string
	h bool
)

func init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.StringVar(&d, "d", "", "Designated `domain` name")
	flag.Usage = usage
}

func main() {
	flag.Parse()

	if d == "" {
		usage()
		os.Exit(1)
	}

	if h {
		flag.Usage()
	}

	//var domain = d + ":443"

	conn, err := net.Dial("tcp", d+":443")
	if err != nil {
		log.Fatal(err)
	}

	client := tls.Client(conn, &tls.Config{
		ServerName: d,
	})
	defer client.Close()

	if err := client.Handshake(); err != nil {
		log.Fatal(err)
	}
	cert := client.ConnectionState().PeerCertificates[0]
	fmt.Println(cert)

	date := time.Now()
	certDate := cert.NotAfter.Local()
	sub := certDate.Sub(date)

	if sub < 168*time.Hour {
		fmt.Println("1")
	}

}

func usage() {
	fmt.Fprintf(os.Stderr, `
Usage: test.exe [-d]

Options:
`)
	flag.PrintDefaults()
}
