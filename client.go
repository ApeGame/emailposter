package main

import (
	"crypto/aes"
	"encoding/hex"
	"fmt"
)

func main() {
	originalText := "kui@ankr.com"

	fmt.Println("Original text:", originalText)

	encrypted, err := encryptAES([]byte(originalText))
	if err != nil {
		panic(err)
	}

	println(encrypted)
}

func encryptAES(data []byte) (string, error) {
	key := "8Q0PljL9dACq2NBCHefyKkHSwWAfHBKx"
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	data = pad(data, blockSize)
	encrypted := make([]byte, len(data))

	for i := 0; i < len(data); i += blockSize {
		block.Encrypt(encrypted[i:i+blockSize], data[i:i+blockSize])
	}

	return hex.EncodeToString(encrypted), nil
}

// pad applies PKCS#7 padding to the given data.
func pad(buf []byte, blockSize int) []byte {
	padLen := blockSize - len(buf)%blockSize
	for i := 0; i < padLen; i++ {
		buf = append(buf, byte(padLen))
	}
	return buf
}
