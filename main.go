package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"pkitool/internal"
	"pkitool/server"
	"pkitool/storage"
	"time"
)

func nextSerialNumber() *big.Int {
	return new(big.Int).SetInt64(time.Now().UnixNano())
}

// returns private key, public DER bytes, cert SN, error
func generateSelfSignedEcdsaCa(name pkix.Name, validity time.Duration) (*internal.EcdsaKeypair, error) {
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

	return &internal.EcdsaKeypair{
		PrivateKey:  priv,
		DerBytes:    der,
		Certificate: cert,
	}, nil
}

func generateEcdsaSubCa(name pkix.Name, validity time.Duration, parentCert *x509.Certificate, parentKey *ecdsa.PrivateKey) (*internal.EcdsaKeypair, error) {
	subPriv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return nil, err
	}

	subTemplate := &x509.Certificate{
		SerialNumber: nextSerialNumber(),
		IsCA:         true,
	}

	der, err := x509.CreateCertificate(rand.Reader, subTemplate, parentCert, &subPriv.PublicKey, parentKey)
	if err != nil {
		return nil, err
	}

	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, err
	}

	return &internal.EcdsaKeypair{
		PrivateKey:  subPriv,
		DerBytes:    der,
		Certificate: cert,
	}, nil
}

func main2() {
	var storer storage.Storer = &storage.FilesystemStorer{BasePath: "test-fs"}

	// Generate CA
	ecdsaCa, err := generateSelfSignedEcdsaCa(pkix.Name{Country: []string{"FR"}}, time.Hour)
	if err != nil {
		panic(err)
	}

	storableCa, err := ecdsaCa.ToStorageKeypair()
	if err != nil {
		panic(err)
	}

	if err := storer.InitPki("test-pki", storableCa); err != nil {
		panic(err)
	}

	// Generate sub 1
	subCa, err := generateEcdsaSubCa(pkix.Name{Country: []string{"FR"}}, time.Hour, ecdsaCa.Certificate, ecdsaCa.PrivateKey)
	if err != nil {
		panic(err)
	}

	storableSubCa, err := subCa.ToStorageKeypair()
	if err != nil {
		panic(err)
	}

	if err := storer.AddSubCA(ecdsaCa.Certificate.SerialNumber.String(), storableSubCa); err != nil {
		panic(err)
	}

	// Generate sub 2
	subCa2, err := generateEcdsaSubCa(pkix.Name{Country: []string{"FR2"}}, time.Hour, subCa.Certificate, subCa.PrivateKey)
	if err != nil {
		panic(err)
	}

	storableSubCa2, err := subCa2.ToStorageKeypair()
	if err != nil {
		panic(err)
	}

	if err := storer.AddSubCA(subCa.Certificate.SerialNumber.String(), storableSubCa2); err != nil {
		panic(err)
	}

	// List PKI
	pki, err := storer.GetFullPki()
	if err != nil {
		panic(err)
	}

	encodedPki, err := json.MarshalIndent(pki, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", encodedPki)

}

func main() {
	router := server.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
