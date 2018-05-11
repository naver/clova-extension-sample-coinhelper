package protocol

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var certMap = make(map[string]x509.Certificate)

func downloadCert(url string) *x509.Certificate {
	if targetCert, found := certMap[url]; found {
		return &targetCert
	}
	tokens := strings.Split(url, "://")
	if tokens[0] != "https" {
		return nil
	}

	response, err := http.Get(url)
	if err != nil {
		log.Println("Error during downloading", url, "-", err)
		return nil
	}
	defer response.Body.Close()
	read, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error during reading")
		return nil
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(read)
	if !ok {
		log.Println("Error during making cert.")
		return nil
	}

	block, _ := pem.Decode(read)
	if block == nil {
		log.Println("Error during making block")
		return nil
	}
	cert, err := x509.ParseCertificates(block.Bytes)
	if err != nil {
		log.Println("Error during parsing certificates")
		return nil
	}
	dnsName := "cek-requests.clova.ai"
	opts := x509.VerifyOptions{
		DNSName: dnsName,
		Roots:   roots,
	}

	if _, err := cert[0].Verify(opts); err != nil {
		log.Println("Failed to verify certificate: " + err.Error())
		return nil
	}

	certMap[url] = *cert[0]
	return cert[0]
}

func CheckSignature(r *http.Request, body []byte) bool {
	signatureStr := r.Header.Get("SignatureCEK")
	signatureDownloadURL := r.Header.Get("SignatureCEKCertChainUrl")

	cert := downloadCert(signatureDownloadURL)
	if cert == nil {
		return false
	}

	hash := sha1.New()
	hash.Write(body)
	hashData := hash.Sum(nil)
	signature, _ := base64.StdEncoding.DecodeString(signatureStr)
	err := rsa.VerifyPKCS1v15(cert.PublicKey.(*rsa.PublicKey), crypto.SHA1, hashData, signature)
	if err != nil {
		return false
	}
	return true
}
