package util

import (
	"crypto/ecdsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"pkitool/models"
)

const (
	CertUrlInject         string = "/cert/%s"
	KeyUrlInject          string = "/key/%s"
	P12UrlInject          string = "/p12/%s"
	SubCaRequestUrlInject string = "/pki/%s/%s/ca"
	CertRequestUrlInject  string = "/pki/%s/%s/cert"
)

func DerToPem(der []byte) (string, error) {
	block := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: der,
	})
	return fmt.Sprintf("%s", block), nil
}

/*
func PemToDer(pem_ string) ([]byte, error) {
	block, _ := pem.Decode([]byte(pem_))
	if block == nil {
		return nil, fmt.Errorf("failed to find a PEM encoded block in the input")
	}
	return block.Bytes, nil
}
*/

func GetDNFromModel(m *models.Dn) pkix.Name {
	return pkix.Name{
		Country:            m.Country,
		Organization:       m.Organization,
		OrganizationalUnit: m.OrgUnit,
		Locality:           m.Locality,
		CommonName:         m.CommonName,
	}
}

func EcdsaKeyToPem(key *ecdsa.PrivateKey) (string, error) {
	keyBytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private key, %s", err)
	}

	block := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	})
	return fmt.Sprintf("%s", block), nil
}
