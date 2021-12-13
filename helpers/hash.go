package helpers

import (
	"crypto"
)

const (
	SHA256 = "SHA256"
	SHA384 = "SHA384"
	SHA512 = "SHA512"
)

func HashBytesUsingAlgorithm(input []byte, signatureAlgorithm string) (crypto.Hash, []byte) {
	alg := hashAlgorithmWithIdentifier(signatureAlgorithm)
	h := alg.New()
	h.Write(input)
	return alg, h.Sum(nil)
}

func hashAlgorithmWithIdentifier(signatureAlgorithm string) crypto.Hash {
	if signatureAlgorithm == SHA256 {
		return crypto.SHA256
	}
	if signatureAlgorithm == SHA384 {
		return crypto.SHA384
	}
	return crypto.SHA512
}
