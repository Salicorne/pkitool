package internal

import (
	"crypto/ecdsa"
	"crypto/x509"
	"pkitool/storage"
	"pkitool/util"
)

type EcdsaKeypair struct {
	PrivateKey  *ecdsa.PrivateKey
	DerBytes    []byte
	Certificate *x509.Certificate
}

func (kp *EcdsaKeypair) ToStorageKeypair() (*storage.KeyPair, error) {
	certPem, err := util.DerToPem(kp.DerBytes)
	if err != nil {
		panic(err)
	}

	keyPem, err := util.EcdsaKeyToPem(kp.PrivateKey)
	if err != nil {
		panic(err)
	}

	return &storage.KeyPair{
		CertPem: certPem,
		CertSN:  kp.Certificate.SerialNumber.String(),
		KeyPem:  keyPem,
	}, nil
}
