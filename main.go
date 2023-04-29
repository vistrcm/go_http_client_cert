package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

// A fundamental concept in `net/http` servers is
// *handlers*. A handler is an object implementing the
// `http.Handler` interface. A common way to write
// a handler is by using the `http.HandlerFunc` adapter
// on functions with the appropriate signature.
func hello(w http.ResponseWriter, req *http.Request) {

	// Functions serving as handlers take a
	// `http.ResponseWriter` and a `http.Request` as
	// arguments. The response writer is used to fill in the
	// HTTP response. Here our simple response is just
	// "hello\n".
	log.Println("got hello")
	fmt.Fprintf(w, "hello w!\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

	// This handler does something a little more
	// sophisticated by reading all the HTTP request
	// headers and echoing them into the response body.
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func createServerConfig(ca, crt, key string) (*tls.Config, error) {
	caCertPEM, err := os.ReadFile(ca)
	if err != nil {
		return nil, err
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(caCertPEM)
	if !ok {
		panic("failed to parse root certificate")
	}

	cert, err := tls.LoadX509KeyPair(crt, key)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		//ClientAuth:            tls.RequireAndVerifyClientCert,
		ClientAuth: tls.RequestClientCert,
		//ClientCAs:             roots,
		VerifyPeerCertificate: verifyPeerCertificate,
	}, nil
}

func verifyPeerCertificate(certs [][]byte, chains [][]*x509.Certificate) error {
	certificates := make([]*x509.Certificate, len(certs))
	for i, asn1Data := range certs {
		cert, err := x509.ParseCertificate(asn1Data)
		if err != nil {
			return errors.New("tls: failed to parse certificate from server: " + err.Error())
		}
		log.Printf("cert: %+v", cert)
		certificates[i] = cert
	}

	log.Printf("chains: %+v", chains)
	for _, c := range chains {
		for _, ci := range c {
			log.Printf("chain: %+v", *ci)
		}
	}

	return nil
}

func main() {

	// We register our handlers on server routes using the
	// `http.HandleFunc` convenience function. It sets up
	// the *default router* in the `net/http` package and
	// takes a function as an argument.
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)

	config, err := createServerConfig("tls/ca.pem", "tls/server_cert.pem", "tls/server_key.pem")
	if err != nil {
		log.Fatal("config failed: %s", err.Error())
	}

	ln, err := tls.Listen("tcp", ":8443", config)
	if err != nil {
		log.Fatal("listen failed: %s", err.Error())
	}

	err = http.Serve(ln, nil)
	if err != nil {
		log.Fatal("serve failed: %s", err.Error())
	}
}
