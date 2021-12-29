package internal

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
)

func nextSerialNumber() *big.Int {
	return new(big.Int).SetInt64(time.Now().UnixNano()) //todo implement thread-safe standalone SN generator
}

func GenerateSelfSignedEcdsaCa(name pkix.Name, validity time.Duration) (*EcdsaKeypair, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return nil, err
	}

	template := &x509.Certificate{
		SerialNumber: nextSerialNumber(),
		Subject:      name,
		Issuer:       name,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(validity),

		KeyUsage: x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		IsCA:     true,
	}

	der, err := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	if err != nil {
		return nil, err
	}

	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, err
	}

	return &EcdsaKeypair{
		PrivateKey:  priv,
		DerBytes:    der,
		Certificate: cert,
	}, nil
}

func GenerateEcdsaSubCa(name pkix.Name, validity time.Duration, parentCert *x509.Certificate, parentKey *ecdsa.PrivateKey) (*EcdsaKeypair, error) {
	subPriv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return nil, err
	}

	subTemplate := &x509.Certificate{
		SerialNumber: nextSerialNumber(),
		IsCA:         true,
		Subject:      name,

		Issuer:    parentCert.Subject,
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(validity),

		KeyUsage: x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
	}

	der, err := x509.CreateCertificate(rand.Reader, subTemplate, parentCert, &subPriv.PublicKey, parentKey)
	if err != nil {
		return nil, err
	}

	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, err
	}

	return &EcdsaKeypair{
		PrivateKey:  subPriv,
		DerBytes:    der,
		Certificate: cert,
	}, nil
}

func GenerateEcdsaKeypair(name pkix.Name, validity time.Duration, parentCert *x509.Certificate, parentKey *ecdsa.PrivateKey) (*EcdsaKeypair, error) {
	subPriv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return nil, err
	}

	template := &x509.Certificate{
		SerialNumber: nextSerialNumber(),
		IsCA:         false,
		Subject:      name,

		Issuer:    parentCert.Subject,
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(validity),
	}

	der, err := x509.CreateCertificate(rand.Reader, template, parentCert, &subPriv.PublicKey, parentKey)
	if err != nil {
		return nil, err
	}

	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, err
	}

	return &EcdsaKeypair{
		PrivateKey:  subPriv,
		DerBytes:    der,
		Certificate: cert,
	}, nil
}
