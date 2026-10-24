// Package ssl generates self-signed RSA server certificates and private keys.
package ssl

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"simonwaldherr.de/go/golibs/atomicfile"
)

const (
	day  = time.Hour * 24
	year = day * 365
)

// Check verifies that both paths exist and name regular files. It does not
// parse the certificate, verify trust or expiration, or check that the keys match.
func Check(certPath string, keyPath string) error {
	for _, path := range []string{certPath, keyPath} {
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("ssl: %s is not a regular file", path)
		}
	}
	return nil
}

// Generate creates a self-signed RSA-2048 server certificate and key. Missing
// certPath, keyPath, countryName, provinceName, organizationName and commonName
// options are prompted for on standard input. Certificates are leaf certificates,
// not certificate authorities; IP common names become IP subject alternative names.
//
// Outputs are prepared in same-directory temporary files before key generation.
// Paths must not alias or differ only by case within the same directory.
// Existing outputs must be regular files, not symlinks. Successful replacements
// use mode 0644 for the public certificate and 0600 for the private key, including
// when an existing key had broader permissions. Each file is replaced by rename;
// replacement of the certificate/key pair is not a single atomic operation.
func Generate(options map[string]string) error {
	scanner := bufio.NewScanner(os.Stdin)
	read := func(name, prompt string) (string, error) {
		if value := options[name]; value != "" {
			return value, nil
		}
		fmt.Print(prompt)
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return "", err
			}
			return "", io.EOF
		}
		return scanner.Text(), nil
	}
	values := make(map[string]string, 6)
	for _, option := range []struct{ name, prompt string }{
		{"certPath", "# Cert filename: "},
		{"keyPath", "# Key filename: "},
		{"countryName", "# Country Name (2 letter code) [AU]: "},
		{"provinceName", "# State or Province Name (full name) [Some-State]: "},
		{"organizationName", "# Organization Name (eg, company) [Lorem Ipsum Inc]: "},
		{"commonName", "# Common Name (eg, YOUR name) [localhost]: "},
	} {
		value, err := read(option.name, option.prompt)
		if err != nil {
			return err
		}
		values[option.name] = value
	}
	certPath, keyPath := values["certPath"], values["keyPath"]
	if certPath == "" || keyPath == "" {
		return fmt.Errorf("ssl: certificate and key require distinct nonempty paths")
	}
	certificate, err := resolveOutput(certPath)
	if err != nil {
		return err
	}
	key, err := resolveOutput(keyPath)
	if err != nil {
		return err
	}
	if strings.EqualFold(filepath.Base(certificate.path), filepath.Base(key.path)) && os.SameFile(certificate.parent, key.parent) ||
		certificate.existing != nil && key.existing != nil && os.SameFile(certificate.existing, key.existing) {
		return fmt.Errorf("ssl: certificate and key require distinct paths")
	}

	err = atomicfile.Write(certificate.path, 0644, func(certOut io.Writer) error {
		return atomicfile.Write(key.path, 0600, func(keyOut io.Writer) error {
			priv, err := rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				return err
			}
			serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
			serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
			if err != nil {
				return err
			}
			notBefore := time.Now().Add(-day)
			template := x509.Certificate{
				SerialNumber: serialNumber,
				Subject: pkix.Name{
					Country:      []string{values["countryName"]},
					Organization: []string{values["organizationName"]},
					Province:     []string{values["provinceName"]},
					CommonName:   values["commonName"],
				},
				NotBefore:             notBefore,
				NotAfter:              notBefore.Add(year * 2),
				KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
				ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
				BasicConstraintsValid: true,
			}
			if ip := net.ParseIP(values["commonName"]); ip != nil {
				template.IPAddresses = []net.IP{ip}
			} else {
				template.DNSNames = []string{values["commonName"]}
			}
			derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
			if err != nil {
				return err
			}
			if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
				return err
			}
			return pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
		})
	})
	if err != nil {
		return err
	}
	log.Printf("written %v\n", certPath)
	log.Printf("written %v\n", keyPath)
	return nil
}

type outputPath struct {
	path     string
	parent   os.FileInfo
	existing os.FileInfo
}

func resolveOutput(path string) (outputPath, error) {
	directory, base := filepath.Split(path)
	if base == "" || base == "." || base == ".." {
		return outputPath{}, &os.PathError{Op: "ssl", Path: path, Err: os.ErrInvalid}
	}
	if directory == "" {
		directory = "."
	}
	// Resolve the raw parent before Abs cleans it: link/.. may name a different
	// directory from its lexical parent. Keep the final leaf unresolved.
	directory, err := filepath.EvalSymlinks(directory)
	if err != nil {
		return outputPath{}, err
	}
	directory, err = filepath.Abs(directory)
	if err != nil {
		return outputPath{}, err
	}
	output := outputPath{path: filepath.Join(directory, base)}
	output.parent, err = os.Stat(directory)
	if err != nil {
		return outputPath{}, err
	}
	output.existing, err = os.Lstat(output.path)
	if err == nil {
		if !output.existing.Mode().IsRegular() {
			return outputPath{}, fmt.Errorf("ssl: %s is not a regular file", path)
		}
	} else if !os.IsNotExist(err) {
		return outputPath{}, err
	}
	return output, nil
}
