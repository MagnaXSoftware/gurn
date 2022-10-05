package main

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)

func main() {
	opts, err := ParseFlags()
	if err != nil {
		fmt.Printf("error parsing options: %s", err.Error())
		return
	}

	crl, err := BuildAndSignCrl(opts)
	if err != nil {
		fmt.Printf("error building CRL: %s", err.Error())
		return
	}

	if nextIdxText, err := big.NewInt(0).Add(opts.Index, big.NewInt(1)).MarshalText(); err != nil {
		fmt.Printf("error saving next index number: %s", err.Error())
		return
	} else if err = os.WriteFile(opts.IndexFile, nextIdxText, 0666); err != nil {
		fmt.Printf("error saving next index number: %s", err.Error())
		return
	}

	outputFile, err := os.Create(opts.Destination)
	if err != nil {
		fmt.Printf("error writting CRL: %s", err.Error())
		return
	}
	block := &pem.Block{
		Bytes: crl,
		Type:  "X509 CRL",
	}
	err = pem.Encode(outputFile, block)
	if err != nil {
		fmt.Printf("error writting CRL: %s", err.Error())
		return
	}
}

func BuildAndSignCrl(opts *Options) ([]byte, error) {
	now := time.Now().UTC()

	revokedList := make([]pkix.RevokedCertificate, 0, len(opts.RevokedCertificates))
	for _, cert := range opts.RevokedCertificates {
		// If the cert would have expired after this update, include it
		if cert.NotAfter.After(now) {
			revokedList = append(revokedList, cert.RevokedCertificate)
		}
	}

	tmpl := &x509.RevocationList{
		Number:              opts.Index,
		ThisUpdate:          now,
		NextUpdate:          now.Add(90 * 24 * time.Hour),
		RevokedCertificates: revokedList,
	}

	return x509.CreateRevocationList(rand.Reader, tmpl, opts.CACert, opts.CAKey)
}
