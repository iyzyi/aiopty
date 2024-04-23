package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"io"
)

var salt = []byte("iyzyi/aiopty")
var sign = []byte("Hello World!")

var errCryptoKey = fmt.Errorf("error crypto key")

func NewCryptoReadWriter(rw io.ReadWriter, key []byte) (crw io.ReadWriter, err error) {
	enc, err := NewCryptoReader(rw, key)
	if err != nil {
		return
	}
	dec, err := NewCryptoWriter(rw, key)
	if err != nil {
		return
	}

	crw = struct {
		io.Reader
		io.Writer
	}{
		Reader: enc,
		Writer: dec,
	}
	return
}

type CryptoReader struct {
	r    io.Reader
	dec  *cipher.StreamReader
	key  []byte
	iv   []byte
	err  error
	init bool
}

func NewCryptoReader(r io.Reader, key []byte) (cr *CryptoReader, err error) {
	// derivation key
	key = pbkdf2.Key(key, salt, 4096, 32, sha256.New)

	cr = &CryptoReader{
		r:   r,
		key: key,
	}
	return
}

func (cr *CryptoReader) Read(b []byte) (n int, err error) {
	if cr.err != nil {
		return 0, cr.err
	}

	if !cr.init {
		// get iv
		iv := make([]byte, aes.BlockSize)
		_, err = io.ReadFull(cr.r, iv)
		if err != nil {
			cr.err = err
			return
		}

		// aes cfb stream
		var block cipher.Block
		block, err = aes.NewCipher(cr.key)
		if err != nil {
			cr.err = err
			return
		}
		cr.dec = &cipher.StreamReader{
			S: cipher.NewCFBDecrypter(block, iv),
			R: cr.r,
		}

		// check sync sign
		_sign := make([]byte, len(sign))
		_, err = io.ReadFull(cr.dec, _sign)
		if err != nil {
			cr.err = err
			return
		}
		if !bytes.Equal(_sign, sign) {
			cr.err = errCryptoKey
			return
		}

		cr.init = true
	}

	n, err = cr.dec.Read(b)
	if err != nil {
		cr.err = err
		return
	}
	return
}

type CryptoWriter struct {
	w    io.Writer
	enc  *cipher.StreamWriter
	key  []byte
	iv   []byte
	err  error
	init bool
}

func NewCryptoWriter(w io.Writer, key []byte) (cw *CryptoWriter, err error) {
	// derivation key
	key = pbkdf2.Key(key, salt, 4096, 32, sha256.New)

	// random iv
	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return
	}

	// aes cfb stream
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	enc := &cipher.StreamWriter{
		S: cipher.NewCFBEncrypter(block, iv),
		W: w,
	}

	cw = &CryptoWriter{
		w:   w,
		enc: enc,
		key: key,
		iv:  iv,
	}
	return
}

func (cw *CryptoWriter) Write(b []byte) (n int, err error) {
	if cw.err != nil {
		return 0, cw.err
	}

	if !cw.init {
		// send iv
		_, err = cw.w.Write(cw.iv)
		if err != nil {
			cw.err = err
			return
		}

		// send sync ciphertext to judge if key is the same.
		_, err = cw.enc.Write(sign)
		if err != nil {
			cw.err = err
			return
		}

		cw.init = true
	}

	n, err = cw.enc.Write(b)
	if err != nil {
		cw.err = err
	}
	return
}
