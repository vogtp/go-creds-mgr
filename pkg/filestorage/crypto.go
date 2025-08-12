package filestorage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

func createHash(key []byte) string {
	hasher := md5.New()
	hasher.Write(key)
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase []byte) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func decrypt(data []byte, passphrase []byte) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

// EncryptedWriter wraps w with an OFB cipher stream.
func EncryptedWriter(key string, w io.Writer) (*cipher.StreamWriter, error) {

	// generate random initial value
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// write clear IV to allow for decryption
	n, err := w.Write(iv)
	if err != nil || n != len(iv) {
		return nil, errors.New("could not write initial value")
	}

	block, err := newBlock(key)
	if err != nil {
		return nil, err
	}

	stream := cipher.NewOFB(block, iv)
	return &cipher.StreamWriter{S: stream, W: w}, nil
}

// EncryptedReader wraps r with an OFB cipher stream.
func EncryptedReader(key string, r io.Reader) (*cipher.StreamReader, error) {

	// read initial value
	iv := make([]byte, aes.BlockSize)
	n, err := r.Read(iv)
	if err != nil || n != len(iv) {
		return nil, errors.New("could not read initial value")
	}

	block, err := newBlock(key)
	if err != nil {
		return nil, err
	}

	stream := cipher.NewOFB(block, iv)
	return &cipher.StreamReader{S: stream, R: r}, nil
}

func newBlock(key string) (cipher.Block, error) {
	hash := md5.Sum([]byte(key))
	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return nil, err
	}
	return block, nil
}
