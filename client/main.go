package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/go-piv/piv-go/piv"
	"net/http"
	"os"
	"strings"
)

func main() {
	cards, err := piv.Cards()
	if err != nil {
		panic(err)
	}

	for _, card := range cards {
		fmt.Printf("Card: %s\n", card)
	}

	// Find a YubiKey and open the reader.
	var yk *piv.YubiKey
	for _, card := range cards {
		if strings.Contains(card, "Yubico YubiKey FIDO+CCID") {
			if yk, err = piv.Open(card); err != nil {
				panic(err)
			}
			break
		}
	}

	fmt.Printf("YubiKey: %s\n", yk)
	if yk == nil {
		panic("nil")
	}

	//// Generate a private key on the YubiKey.
	//key := piv.Key{
	//	Algorithm:   piv.AlgorithmEC256,
	//	PINPolicy:   piv.PINPolicyAlways,
	//	TouchPolicy: piv.TouchPolicyNever,
	//}
	//pub, err := yk.GenerateKey(piv.DefaultManagementKey, piv.SlotAuthentication, key)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("Public key: %s\n", pub)

	cert, err := yk.Certificate(piv.SlotAuthentication)
	if err != nil {
		panic(err)
	}

	auth := piv.KeyAuth{PIN: piv.DefaultPIN}

	priv, err := yk.PrivateKey(piv.SlotAuthentication, cert.PublicKey, auth)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Private key: %s\n", priv)

	//certificate, err := tls.LoadX509KeyPair(
	//	"/Users/vist/proj/github.com/vistrcm/go_http_client_cert/tls/client_cert.pem",
	//	"/Users/vist/proj/github.com/vistrcm/go_http_client_cert/tls/client_key.pem")
	//if err != nil {
	//	panic(err)
	//}

	caCert, _ := os.ReadFile("/Users/vist/proj/github.com/vistrcm/go_http_client_cert/tls/ca.pem")
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			//Certificates:       []tls.Certificate{certificate},
			RootCAs: caCertPool,
			Certificates: []tls.Certificate{
				{
					Certificate: [][]byte{cert.Raw},
					PrivateKey:  priv,
				},
			},
		},
	}

	//s := &http.Server{
	//	Addr: ":8444",
	//	TLSConfig: &tls.Config{
	//		Certificates: []tls.Certificate{
	//			{
	//				Certificate: [][]byte{cert.Raw},
	//				PrivateKey:  priv,
	//			},
	//		},
	//	},
	//	Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	//	}),
	//}
	//
	//log.Fatal(s.ListenAndServeTLS("", ""))

	client := &http.Client{Transport: tr}
	resp, err := client.Get("https://localhost:8443/hello")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response: %+v\n", resp)
}
