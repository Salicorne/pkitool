package internal

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
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

func EcdsaFromStorageKeypair(kp *storage.KeyPair) (*EcdsaKeypair, error) {
	keyBlock, _ := pem.Decode([]byte(kp.KeyPem))
	if keyBlock == nil {
		return nil, fmt.Errorf("failed to load key PEM from storage")
	}

	privKey, err := x509.ParseECPrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse EC key: %s", err.Error())
	}

	certBlock, _ := pem.Decode([]byte(kp.CertPem))
	if keyBlock == nil {
		return nil, fmt.Errorf("failed to load certificate PEM from storage")
	}

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %s", err.Error())
	}

	return &EcdsaKeypair{
		PrivateKey:  privKey,
		Certificate: cert,
		DerBytes:    certBlock.Bytes,
	}, nil
}
