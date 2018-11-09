// Written by Nadim Kobeissi, <nadim@symbolic.software> November 2018
// Copyright (c) 2018 Nadim Kobeissi
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"crypto/rand"
	"errors"

	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/curve25519"
)

type dgKeypair struct {
	sk [32]byte
	pk [32]byte
}

type dgCryptojar struct {
	ciphertext []byte
	nonce      [24]byte
}

func dgCryptoEdhCheckSharedSecret(ss [32]byte) bool {
	good := false
	for _, v := range ss[:32] {
		// TODO: Try to approach constant time before using
		// in the rest of the code.
		if v > 0 {
			good = true
		}
	}
	if !good {
		dgErrorCritical(errors.New("bad shared secret"))
	}
	return good
}

func dgCryptoEdhGenerate() dgKeypair {
	var pk [32]byte
	var sk [32]byte
	_, err := rand.Read(sk[:32])
	dgErrorCritical(err)
	curve25519.ScalarBaseMult(&pk, &sk)
	return dgKeypair{sk, pk}
}

func dgCryptoEdhSharedSecret(sk [32]byte, pk [32]byte) [32]byte {
	var ss [32]byte
	curve25519.ScalarMult(&ss, &sk, &pk)
	dgCryptoEdhCheckSharedSecret(ss)
	return ss
}

func dgCryptoEncFile(key [32]byte, file []byte) dgCryptojar {
	var nonce [24]byte
	enc, err := chacha20poly1305.NewX(key[:32])
	dgErrorCritical(err)
	_, err = rand.Read(nonce[:24])
	dgErrorCritical(err)
	ciphertext := enc.Seal(nil, nonce[:24], file, nil)
	return dgCryptojar{ciphertext, nonce}
}

func dgCryptoDecFile(key [32]byte, c dgCryptojar) []byte {
	enc, err := chacha20poly1305.NewX(key[:32])
	dgErrorCritical(err)
	plaintext, err := enc.Open(nil, c.nonce[:24], c.ciphertext, nil)
	dgErrorCritical(err)
	return plaintext
}
