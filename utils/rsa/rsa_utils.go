package utils

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"log"
)

const (
	rsaKeySize = 2048
)

// GenerateKeyPair Copy from https://gist.github.com/miguelmota/3ea9286bd1d3c2a985b67cac4ba2130a
// GenerateKeyPair generates a new key pair
func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		log.Fatal(err)
	}
	return privkey, &privkey.PublicKey
}

// PrivateKeyToBytes private key to bytes
func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return privBytes
}

// PublicKeyToBytes public key to bytes
func PublicKeyToBytes(pub *rsa.PublicKey) []byte {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		log.Fatal(err)
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes
}

// BytesToPrivateKey bytes to private key
func BytesToPrivateKey(priv []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		log.Fatal(err)
	}
	return key
}

// BytesToPublicKey bytes to public key
func BytesToPublicKey(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		log.Fatal(err)
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		log.Fatal("not ok")
	}
	return key
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) []byte {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		log.Fatal(err)
	}
	return ciphertext
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) []byte {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		log.Fatal(err)
	}
	return plaintext
}

// Copy from https://gist.github.com/wongoo/2b974a9594627114bea3e53c794980cd
// hash
func hash(data []byte) []byte {
	s := sha1.Sum(data)
	return s[:]

}

// GenerateKey Copy from https://gist.github.com/wongoo/2b974a9594627114bea3e53c794980cd
// GenerateKey
func GenerateKey() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	pri, err := rsa.GenerateKey(rand.Reader, rsaKeySize)
	if err != nil {
		return nil, nil, err
	}
	return pri, &pri.PublicKey, nil
}

func GenerateKeyBytes() (privateBytes, publicBytes []byte, err error) {
	pri, pub, err := GenerateKey()
	if err != nil {
		return nil, nil, err
	}
	priBytes, err := x509.MarshalPKCS8PrivateKey(pri)
	if err != nil {
		return nil, nil, err
	}
	pubBytes := x509.MarshalPKCS1PublicKey(pub)
	return priBytes, pubBytes, nil
}

func GenerateKey64() (pri64, pub64 string, err error) {
	pri, pub, err := GenerateKeyBytes()
	if err != nil {
		return "", "", nil
	}
	return base64.StdEncoding.EncodeToString(pri),
		base64.StdEncoding.EncodeToString(pub),
		nil
}

func PublicKeyFrom(key []byte) (*rsa.PublicKey, error) {
	pubInterface, err := x509.ParsePKIXPublicKey(key)
	if err != nil {
		return nil, err
	}
	pub, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		log.Fatal(ok)
	}
	return pub, nil
}

func PublicKeyFrom64(key string) (*rsa.PublicKey, error) {
	b, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	return PublicKeyFrom(b)
}

func PrivateKeyFrom(key []byte) (*rsa.PrivateKey, error) {
	pri, err := x509.ParsePKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}
	p, ok := pri.(*rsa.PrivateKey)
	if !ok {
		log.Fatal(ok)
	}
	return p, nil
}

func PrivateKeyFrom64(key string) (*rsa.PrivateKey, error) {
	b, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	return PrivateKeyFrom(b)
}

func PublicEncrypt(key *rsa.PublicKey, data []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, key, data)
}

func PublicSign(key *rsa.PublicKey, data []byte) ([]byte, error) {
	return PublicEncrypt(key, hash(data))
}

func PublicVerify(key *rsa.PublicKey, sign, data []byte) error {
	return rsa.VerifyPKCS1v15(key, crypto.SHA1, hash(data), sign)
}

func PrivateDecrypt(key *rsa.PrivateKey, data []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, key, data)
}

func PrivateSign(key *rsa.PrivateKey, data []byte) ([]byte, error) {
	return rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA1, hash(data))
}

func PrivateVerify(key *rsa.PrivateKey, sign, data []byte) error {
	h, err := PrivateDecrypt(key, sign)
	if err != nil {
		return err
	}
	if !bytes.Equal(h, hash(data)) {
		return rsa.ErrVerification
	}
	return nil
}
