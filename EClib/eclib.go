package eclib

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"math/big"
	"os"
	"strings"

	"github.com/gopherjs/gopherjs/js"
)

// EC struct is a struct full of the elliptic curve methods
type EC struct {
	ECGenerateKey func() (string, string)
	ECHashString  func(string) string
	ECSign        func(string, string, string) string
	ECVerify      func(string, string, string) bool
}

func main() {
	ec := EC{
		ECGenerateKey: ECGenerateKey,
		ECHashString:  ECHashString,
		ECSign:        ECSign,
		ECVerify:      ECVerify,
	}
	js.Global.Set("ec", ec)
}

// ECGenerateKey is a function which return a private key and a public key
func ECGenerateKey() (string, string) {
	pubkeyCurve := elliptic.P256() //see http://golang.org/pkg/crypto/elliptic/#P256

	privatekey := new(ecdsa.PrivateKey)
	privatekey, _ = ecdsa.GenerateKey(pubkeyCurve, rand.Reader) // this generates a public & private key pair

	var pubkey ecdsa.PublicKey
	pubkey = privatekey.PublicKey
	pkeyBytes := elliptic.Marshal(pubkey, pubkey.X, pubkey.Y)
	encoded := base64.StdEncoding.EncodeToString(pkeyBytes)
	return privatekey.D.String(), encoded
}

// ECHashString hashes a string with the sha256 algorithm
func ECHashString(toHash string) string {
	var h hash.Hash
	h = sha256.New()

	io.WriteString(h, toHash)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// ECSign signs a string with a private and public key and returns the sign string
func ECSign(toSign string, D string, publicKey string) string {
	var d big.Int
	d.SetString(D, 10)
	pkey, _ := base64.StdEncoding.DecodeString(publicKey)
	toSignBytes, err := base64.StdEncoding.DecodeString(toSign)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	x, y := elliptic.Unmarshal(elliptic.P256(), []byte(pkey))
	privateKey := ecdsa.PrivateKey{PublicKey: ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}, D: &d}
	r, s, serr := ecdsa.Sign(rand.Reader, &privateKey, toSignBytes)
	if serr != nil {
		fmt.Println(serr)
		os.Exit(1)
	}

	signature := string(r.String()) + "-" + string(s.String())

	return signature
}

// ECVerify checks if a string is signed with a certain public key
func ECVerify(toVerify string, signature string, pubkey string) bool {
	pubkeyBytes, _ := base64.StdEncoding.DecodeString(pubkey)
	x, y := elliptic.Unmarshal(elliptic.P256(), pubkeyBytes)
	publicKey := ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}
	rs := strings.Split(signature, "-")
	if len(rs) == 2 {
		var r, s big.Int
		r.SetString(rs[0], 10)
		s.SetString(rs[1], 10)
		toVerifyBytes, _ := base64.StdEncoding.DecodeString(toVerify)
		return ecdsa.Verify(&publicKey, toVerifyBytes, &r, &s)
	}

	return false
}
