package concordium

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
)

const signKeySize = 32

type SignKey interface {
	Decode() ([]byte, error)
}

type Credentials map[uint8]map[uint8]KeyPair

type KeyPair struct {
	SignKey   SignKey
	VerifyKey string
}

type DecryptedSignKey string

func (k DecryptedSignKey) Decode() ([]byte, error) {
	return hex.DecodeString(string(k))
}

type EncryptedSignKey struct {
	Password   string                   `json:"password"`
	Metadata   EncryptedSignKeyMetadata `json:"metadata"`
	CipherText string                   `json:"cipherText"`
}

type EncryptedSignKeyMetadata struct {
	Iterations           int    `json:"iterations"`
	Salt                 string `json:"salt"`
	InitializationVector string `json:"initializationVector"`
}

func (k *EncryptedSignKey) Decode() ([]byte, error) {
	if k.Password == "" {
		return nil, fmt.Errorf("empty password")
	}
	salt, err := base64.StdEncoding.DecodeString(k.Metadata.Salt)
	if err != nil {
		return nil, fmt.Errorf("salt decode: %w", err)
	}
	iv, err := base64.StdEncoding.DecodeString(k.Metadata.InitializationVector)
	if err != nil {
		return nil, fmt.Errorf("initialization vector decode: %w", err)
	}
	block, err := aes.NewCipher(
		pbkdf2.Key([]byte(k.Password), salt, k.Metadata.Iterations, signKeySize, sha256.New),
	)
	if err != nil {
		return nil, fmt.Errorf("cipher create: %w", err)
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	cipherText, err := base64.StdEncoding.DecodeString(k.CipherText)
	if err != nil {
		return nil, fmt.Errorf("cipher text decode: %w", err)
	}
	dec := make([]byte, len(cipherText))
	mode.CryptBlocks(dec, cipherText)
	key := make([]byte, signKeySize)
	_, err = hex.Decode(key, dec[1:65])
	return key, err
}
