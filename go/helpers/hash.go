package helpers

import (
	"crypto"
)

func HashBytesUsingAlgorithm(input []byte, signatureAlgorithm string) (crypto.Hash, []byte) {
	alg := hashAlgorithmWithIdentifier(signatureAlgorithm)
	h := alg.New()
	h.Write(input)
	return alg, h.Sum(nil)
}

func hashAlgorithmWithIdentifier(signatureAlgorithm string) crypto.Hash {
	if signatureAlgorithm == "SHA256" {
		return crypto.SHA256
	}
	return crypto.SHA512
}
