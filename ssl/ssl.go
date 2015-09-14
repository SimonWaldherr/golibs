// Generate SSL Certs via this simple to use ssl package.
package ssl

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"
)

const (
	day  = time.Hour * 24
	year = day * 365
)

func Check(certPath string, keyPath string) error {
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return err
	} else if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return err
	}
	return nil
}

func Generate(options map[string]string) error {
	var err error
	var certOut *os.File
	var derBytes []byte

	priv, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		log.Printf("failed to generate private key: %s", err)
		return err
	}

	notBefore := time.Now().Add(day * -1)
	notAfter := notBefore.Add(year * 2)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Printf("failed to generate serial number: %s", err)
		return err
	}

	scanner := bufio.NewScanner(os.Stdin)

	var certPath string
	if options["certPath"] == "" {
		fmt.Print("# Cert filename: ")
		scanner.Scan()
		certPath = scanner.Text()
	} else {
		certPath = options["certPath"]
	}

	var keyPath string
	if options["keyPath"] == "" {
		fmt.Print("# Key filename: ")
		scanner.Scan()
		keyPath = scanner.Text()
	} else {
		keyPath = options["keyPath"]
	}

	var countryName string
	if options["countryName"] == "" {
		fmt.Print("# Country Name (2 letter code) [AU]: ")
		scanner.Scan()
		countryName = scanner.Text()
	} else {
		countryName = options["countryName"]
	}

	var provinceName string
	if options["provinceName"] == "" {
		fmt.Print("# State or Province Name (full name) [Some-State]: ")
		scanner.Scan()
		provinceName = scanner.Text()
	} else {
		provinceName = options["provinceName"]
	}

	var organizationName string
	if options["organizationName"] == "" {
		fmt.Print("# Organization Name (eg, company) [Lorem Ipsum Inc]: ")
		scanner.Scan()
		organizationName = scanner.Text()
	} else {
		organizationName = options["organizationName"]
	}

	var commonName string
	if options["commonName"] == "" {
		fmt.Print("# Common Name (eg, YOUR name) [localhost]: ")
		scanner.Scan()
		commonName = scanner.Text()
	} else {
		commonName = options["commonName"]
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:      []string{string(countryName)},
			Organization: []string{string(organizationName)},
			Province:     []string{string(provinceName)},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{string(commonName)},
		IsCA:                  true,
	}

	if derBytes, err = x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv); err != nil {
		log.Printf("Failed to create certificate: %s", err)
		return err
	}

	if certOut, err = os.Create(certPath); err != nil {
		log.Printf("failed to open "+certPath+" for writing: %s", err)
		return err
	}

	if err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return err
	}
	if err = certOut.Close(); err != nil {
		return err
	}

	log.Printf("written %v\n", certPath)

	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Print("failed to open "+keyPath+" for writing:", err)
		return err
	}

	if err = pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		return err
	}

	if err = keyOut.Close(); err != nil {
		return err
	}
	log.Printf("written %v\n", keyPath)
	return nil
}
