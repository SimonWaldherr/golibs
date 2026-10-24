package ssl

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func certificateOptions(directory, name string) map[string]string {
	return map[string]string{
		"certPath":    filepath.Join(directory, "server.cert"),
		"keyPath":     filepath.Join(directory, "server.key"),
		"countryName": "DE", "provinceName": "Bavaria", "organizationName": "Example",
		"commonName": name,
	}
}

func TestGenerateServerCertificateAndPrivateReplacement(t *testing.T) {
	for _, name := range []string{"example.test", "127.0.0.1", "::1"} {
		t.Run(name, func(t *testing.T) {
			options := certificateOptions(t.TempDir(), name)
			// Replacing an old permissive key must not retain its old mode.
			if err := os.WriteFile(options["keyPath"], []byte("old private key"), 0644); err != nil {
				t.Fatal(err)
			}
			if err := os.Chmod(options["keyPath"], 0644); err != nil {
				t.Fatal(err)
			}
			if err := Generate(options); err != nil {
				t.Fatal(err)
			}
			if err := Check(options["certPath"], options["keyPath"]); err != nil {
				t.Fatal(err)
			}
			certificate, err := os.ReadFile(options["certPath"])
			if err != nil {
				t.Fatal(err)
			}
			key, err := os.ReadFile(options["keyPath"])
			if err != nil {
				t.Fatal(err)
			}
			if _, err := tls.X509KeyPair(certificate, key); err != nil {
				t.Fatalf("certificate and key do not match: %v", err)
			}
			block, _ := pem.Decode(certificate)
			if block == nil || block.Type != "CERTIFICATE" {
				t.Fatal("certificate output is not PEM")
			}
			parsed, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				t.Fatal(err)
			}
			if parsed.IsCA || parsed.KeyUsage&x509.KeyUsageCertSign != 0 {
				t.Fatal("server certificate was given certificate-authority privileges")
			}
			if err := parsed.VerifyHostname(name); err != nil {
				t.Fatalf("SAN is not usable for %q: %v", name, err)
			}
			if len(parsed.ExtKeyUsage) != 1 || parsed.ExtKeyUsage[0] != x509.ExtKeyUsageServerAuth {
				t.Fatalf("unexpected extended key usages: %v", parsed.ExtKeyUsage)
			}
			keyBlock, _ := pem.Decode(key)
			if keyBlock == nil || keyBlock.Type != "RSA PRIVATE KEY" {
				t.Fatal("RSA PEM compatibility changed")
			}
			private, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
			if err != nil || private.N.BitLen() != 2048 {
				t.Fatalf("invalid RSA-2048 private key: %v", err)
			}
			info, err := os.Stat(options["keyPath"])
			if err != nil {
				t.Fatal(err)
			}
			if runtime.GOOS != "windows" && info.Mode().Perm() != 0600 {
				t.Fatalf("replacement key permissions = %o, want 600", info.Mode().Perm())
			}
			entries, err := os.ReadDir(filepath.Dir(options["certPath"]))
			if err != nil || len(entries) != 2 {
				t.Fatalf("temporary PEM files were left behind: %v, %v", entries, err)
			}
		})
	}
}

func TestGenerateRejectsInvalidOutputsWithoutReplacingCertificate(t *testing.T) {
	for _, failure := range []string{"missing key parent", "key directory", "certificate directory", "same output path"} {
		t.Run(failure, func(t *testing.T) {
			directory := t.TempDir()
			options := certificateOptions(directory, "example.test")
			if err := os.WriteFile(options["certPath"], []byte("keep certificate"), 0644); err != nil {
				t.Fatal(err)
			}
			switch failure {
			case "missing key parent":
				options["keyPath"] = filepath.Join(directory, "missing", "server.key")
			case "key directory":
				options["keyPath"] = directory
			case "certificate directory":
				options["certPath"] = directory
			case "same output path":
				options["keyPath"] = options["certPath"]
			}
			if err := Generate(options); err == nil {
				t.Fatal("invalid output path was accepted")
			}
			certificate, err := os.ReadFile(filepath.Join(directory, "server.cert"))
			if err != nil || string(certificate) != "keep certificate" {
				t.Fatalf("failed generation changed existing certificate: %q, %v", certificate, err)
			}
			entries, err := os.ReadDir(directory)
			if err != nil || len(entries) != 1 {
				t.Fatalf("failed generation left temporary files: %v, %v", entries, err)
			}
		})
	}
}

func TestGenerateRejectsSymlinkOutputs(t *testing.T) {
	directory := t.TempDir()
	options := certificateOptions(directory, "example.test")
	target := filepath.Join(directory, "unrelated.key")
	if err := os.WriteFile(target, []byte("do not replace"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, options["keyPath"]); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if err := Generate(options); err == nil {
		t.Fatal("symlink key output was accepted")
	}
	content, err := os.ReadFile(target)
	if err != nil || string(content) != "do not replace" {
		t.Fatalf("symlink target changed: %q, %v", content, err)
	}
	if _, err := os.Stat(options["certPath"]); !os.IsNotExist(err) {
		t.Fatalf("failed generation published a certificate: %v", err)
	}
}

func TestGenerateRejectsAliasedOutputs(t *testing.T) {
	directory := t.TempDir()
	options := certificateOptions(directory, "example.test")
	if err := os.WriteFile(options["certPath"], []byte("keep file"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Link(options["certPath"], options["keyPath"]); err != nil {
		t.Skipf("hard links unavailable: %v", err)
	}
	if err := Generate(options); err == nil {
		t.Fatal("certificate and key paths aliasing the same file were accepted")
	}
	for _, path := range []string{options["certPath"], options["keyPath"]} {
		content, err := os.ReadFile(path)
		if err != nil || string(content) != "keep file" {
			t.Fatalf("aliased output was changed: %q, %v", content, err)
		}
	}
}

func TestGenerateRejectsCaseOnlyMissingOutputs(t *testing.T) {
	directory := t.TempDir()
	options := certificateOptions(directory, "example.test")
	options["keyPath"] = filepath.Join(directory, "SERVER.CERT")
	if err := Generate(options); err == nil {
		t.Fatal("case-only output names were accepted")
	}
	entries, err := os.ReadDir(directory)
	if err != nil || len(entries) != 0 {
		t.Fatalf("rejected output names created files: %v, %v", entries, err)
	}
}

func TestGeneratePreservesSymlinkParentSemantics(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symbolic links may require elevated privileges")
	}
	for _, sameTarget := range []bool{false, true} {
		t.Run(map[bool]string{false: "distinct targets", true: "same target"}[sameTarget], func(t *testing.T) {
			directory := t.TempDir()
			lexical := filepath.Join(directory, "lexical")
			actual := filepath.Join(directory, "actual")
			child := filepath.Join(actual, "child")
			for _, parent := range []string{lexical, child} {
				if err := os.MkdirAll(parent, 0700); err != nil {
					t.Fatal(err)
				}
			}
			link := filepath.Join(lexical, "link")
			if err := os.Symlink(child, link); err != nil {
				t.Fatal(err)
			}
			options := certificateOptions(directory, "example.test")
			// Keep link/.. intact: its actual parent is not its lexical parent.
			options["certPath"] = link + string(os.PathSeparator) + ".." + string(os.PathSeparator) + "same"
			options["keyPath"] = filepath.Join(lexical, "same")
			if sameTarget {
				options["keyPath"] = filepath.Join(actual, "same")
				if err := Generate(options); err == nil {
					t.Fatal("aliased missing targets were accepted")
				}
			} else {
				if err := Generate(options); err != nil {
					t.Fatal(err)
				}
				cert, err := os.ReadFile(filepath.Join(actual, "same"))
				if err != nil {
					t.Fatal(err)
				}
				key, err := os.ReadFile(filepath.Join(lexical, "same"))
				if err != nil {
					t.Fatal(err)
				}
				if _, err := tls.X509KeyPair(cert, key); err != nil {
					t.Fatalf("outputs do not form a key pair: %v", err)
				}
			}
			for _, parent := range []string{lexical, actual} {
				entries, err := os.ReadDir(parent)
				if err != nil {
					t.Fatal(err)
				}
				for _, entry := range entries {
					if strings.HasPrefix(entry.Name(), ".atomicfile-") {
						t.Fatalf("temporary output remained: %s", entry.Name())
					}
				}
			}
		})
	}
}

func TestCheckRejectsMissingFilesAndDirectories(t *testing.T) {
	directory := t.TempDir()
	regular := filepath.Join(directory, "file")
	if err := os.WriteFile(regular, []byte("placeholder"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, paths := range [][2]string{
		{filepath.Join(directory, "missing"), regular},
		{regular, filepath.Join(directory, "missing")},
		{directory, regular}, {regular, directory},
		{filepath.Join(regular, "child"), regular},
	} {
		if err := Check(paths[0], paths[1]); err == nil {
			t.Errorf("Check accepted invalid paths %v", paths)
		}
	}
	if err := Check(regular, regular); err != nil {
		t.Fatalf("Check must only test regular-file existence: %v", err)
	}
}

func TestGenerateInputEOF(t *testing.T) {
	input, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	defer input.Close()
	previous := os.Stdin
	os.Stdin = input
	defer func() { os.Stdin = previous }()
	if err := Generate(nil); !errors.Is(err, io.EOF) {
		t.Fatalf("Generate with exhausted input returned %v, want EOF", err)
	}
}
