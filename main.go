package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/dchest/uniuri"
	"github.com/yawning/bulb"
	"github.com/yawning/bulb/utils/pkcs1"
)

const (
	controlPath    = "/run/tor/control"
	authCookiePath = "/run/tor/control.authcookie"
	rsaKeySize     = 1024
	nameLength     = 22 // provides ~128 bits of security
)

var ignoreExts = map[string]bool{
	".bz2": true,
	".gz":  true,
	".xz":  true,
}

func generateFilename(source string) string {
	barename := uniuri.NewLen(nameLength)
	ext := filepath.Ext(source)
	if ignoreExts[ext] {
		ext = filepath.Ext(source[:len(source)-len(ext)]) + ext
	}

	return barename + ext
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if path.Clean(r.URL.Path) != "/" {
		http.NotFound(w, r)
		return
	}

	io.WriteString(w, "I’m just a happy little web server.\n")
}

func makeFileHandler(path string) func(http.ResponseWriter, *http.Request) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Frame-Options", "DENY")
		w.Header().Add("Content-Security-Policy", "default-src 'none'; referrer none;")

		f, err := os.Open(path)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		http.ServeFile(w, r, path)
	}
	return handler
}

func main() {
	c, err := bulb.Dial("unix", controlPath)
	if err != nil {
		log.Fatal("Failed to connect to control port: ", err)
	}
	defer c.Close()

	// read auth cookie and authenticate with it
	cookie, err := ioutil.ReadFile(authCookiePath)
	if err != nil {
		log.Fatal("Failed to read auth cookie: ", err)
	}

	if err := c.Authenticate(hex.EncodeToString(cookie)); err != nil {
		log.Fatal("Authentication failed: ", err)
	}

	pk, err := rsa.GenerateKey(rand.Reader, rsaKeySize)
	if err != nil {
		log.Fatal("Failed to generate RSA key: ", err)
	}

	addr, err := pkcs1.OnionAddr(&pk.PublicKey)
	if err != nil {
		log.Fatal("Failed to derive onion address: ", err)
	}
	log.Printf("%v", addr)

	// Note: this requires Tor 0.2.7.x
	l, err := c.Listener(80, pk)
	if err != nil {
		log.Fatal("Failed to get listener: ", err)
	}
	defer l.Close()

	log.Printf("Listener: %s", l.Addr())

	http.HandleFunc("/", indexHandler)

	flag.Parse()
	for _, arg := range flag.Args() {
		path := fmt.Sprintf("/%v", generateFilename(arg))
		log.Printf("%v: %v", arg, path)
		http.HandleFunc(path, makeFileHandler(arg))
	}

	http.Serve(l, nil)
}
