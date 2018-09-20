package protocol

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var publicKey *rsa.PublicKey

const publicKeyDownloadURL = "https://clova.ai/.well-known/signature-public-key.pem"

func downloadPublicKey() bool {
	tokens := strings.Split(publicKeyDownloadURL, "://")
	if tokens[0] != "https" {
		return false
	}

	response, err := http.Get(publicKeyDownloadURL)
	if err != nil {
		log.Println("Error during downloading", publicKeyDownloadURL, "-", err)
		return false
	}
	defer response.Body.Close()
	read, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error during reading")
		return false
	}

	block, _ := pem.Decode([]byte(read))
	downloadedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	publicKey = downloadedKey.(*rsa.PublicKey)
	log.Println("Download public key complete")
	return true
}

func CheckSignature(r *http.Request, body []byte) bool {
	signatureStr := r.Header.Get("SignatureCEK")

	if publicKey == nil && !downloadPublicKey() {
		return false
	}

	hash := crypto.SHA256.New()
	hash.Write(body)
	hashData := hash.Sum(nil)
	signature, _ := base64.StdEncoding.DecodeString(signatureStr)
	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashData, signature)
	if err != nil {
		return false
	}
	return true
}
