package main

import (
	"crypto"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"math/big"
	"os"
	"time"
)

type Options struct {
	CACert      *x509.Certificate
	CAKey       crypto.Signer
	Index       *big.Int
	IndexFile   string
	Destination string

	RevokedCertificates []RevokedCertificate
}

type RevokedCertificate struct {
	pkix.RevokedCertificate

	NotAfter time.Time
}

var (
	caCertPath         string
	caKeyPath          string
	crlIndexPath       string
	crlDestinationPath string
)

func init() {
	flag.StringVar(&caCertPath, "ca-cert", "", "Path to the CA Certificate")
	flag.StringVar(&caKeyPath, "ca-key", "", "Path to the CA Key")
	flag.StringVar(&crlIndexPath, "index", "crlnumber", "Path to the file containing the next CRL serial number")
	flag.StringVar(&crlDestinationPath, "destination", "", "Path of the generated CRL")
}

func ParseFlags() (*Options, error) {
	flag.Parse()

	opt := &Options{}

	if "" == caCertPath {
		return nil, fmt.Errorf("ca-cert path cannot be empty")
	}
	raw, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, err
	}
	block, rest := pem.Decode(raw)
	for {
		if block == nil || len(rest) == 0 {
			if opt.CACert == nil {
				return nil, fmt.Errorf("no PEM data in %s", caCertPath)
			} else {
				break
			}
		}
		if block.Type == "CERTIFICATE" {
			opt.CACert, err = x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, err
			}
		}
		if signer, err := parsePublicKeyBlock(block); err != nil {
			return nil, err
		} else if signer != nil {
			opt.CAKey = signer
		}

		block, rest = pem.Decode(rest)
	}

	if opt.CAKey == nil {
		// We don't have a key, so read the ca-key path
		if "" == caKeyPath {
			return nil, fmt.Errorf("ca-key path cannot be empty if the key isn't combined in the cert file")
		}
		raw, err := os.ReadFile(caKeyPath)
		if err != nil {
			return nil, err
		}
		block, rest := pem.Decode(raw)
		for {
			if block == nil || len(rest) == 0 {
				if opt.CAKey == nil {
					return nil, fmt.Errorf("no PEM data in %s", caKeyPath)
				} else {
					break
				}
			}
			if signer, err := parsePublicKeyBlock(block); err != nil {
				return nil, err
			} else if signer != nil {
				opt.CAKey = signer
			}

			block, rest = pem.Decode(rest)
		}
	}

	if "" == crlIndexPath {
		return nil, fmt.Errorf("index cannot be empty")
	}
	raw, err = os.ReadFile(crlIndexPath)
	if err != nil {
		return nil, err
	}
	opt.Index = &big.Int{}
	err = opt.Index.UnmarshalText(raw)
	if err != nil {
		return nil, err
	}
	opt.IndexFile = crlIndexPath

	if "" == crlDestinationPath {
		return nil, fmt.Errorf("destination cannot be empty")
	}
	opt.Destination = crlDestinationPath

	return opt, nil
}

func parsePublicKeyBlock(block *pem.Block) (crypto.Signer, error) {
	if block.Type == "PUBLIC KEY" {
		// We have a PKCS8 cert
		akey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		if key, ok := akey.(crypto.Signer); ok {
			return key, nil
		} else {
			return nil, fmt.Errorf("unable to parse private key from %s", caCertPath)
		}
	}
	if block.Type == "RSA PRIVATE KEY" {
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	}
	if block.Type == "EC PRIVATE KEY" {
		return x509.ParseECPrivateKey(block.Bytes)
	}
	return nil, nil
}
