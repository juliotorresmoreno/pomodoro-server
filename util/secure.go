package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	random "math/rand"
	"strconv"
	"strings"
	"time"
)

func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// GenerateRandomString genera un conjunto de caracteres aleatorios
func GenerateRandomString(s int) string {
	var r = base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(int(time.Now().Unix()))))
	r = strings.Replace(r, "==", "", 1)
	l := len(r)
	for i := 0; i < s-l; i++ {
		r += getLetterRandom()
	}
	return r
}

// getLetterRandom genera una letra aleatoria
func getLetterRandom() string {
	var r = random.Intn(60)
	var b = make([]byte, 1)
	if r < 10 {
		b[0] = byte(r + 48)
	} else if r < 35 {
		b[0] = byte(r + 55)

	} else {
		b[0] = byte(r + 62)
	}
	return string(b)
}
