package storage

import (
	"fmt"
	"pkitool/models"
	"pkitool/util"
)

type KeyPair struct {
	CertPem string
	CertSN  string
	KeyPem  string
}

type PkiNode struct {
	SerialNumber string
	Ca           bool
	Children     []*PkiNode
}

func (n *PkiNode) ToPkiDetails(pkiName string) *models.PkiDetails {
	children := []models.PkiDetails{}

	for c := range n.Children {
		children = append(children, *n.Children[c].ToPkiDetails(pkiName))
	}

	return &models.PkiDetails{
		SerialNumber:    n.SerialNumber,
		Ca:              n.Ca,
		CertUrl:         fmt.Sprintf(util.CertUrlInject, n.SerialNumber),
		KeyUrl:          fmt.Sprintf(util.KeyUrlInject, n.SerialNumber),
		P12Url:          fmt.Sprintf(util.P12UrlInject, n.SerialNumber),
		SubCaRequestUrl: fmt.Sprintf(util.SubCaRequestUrlInject, pkiName, n.SerialNumber),
		CertRequestUrl:  fmt.Sprintf(util.CertRequestUrlInject, pkiName, n.SerialNumber),
		Children:        children,
	}
}

// Main interface implemented by storage engines
type Storer interface {
	InitPki(name string, keypair *KeyPair) error
	AddSubCA(parentSN string, keypair *KeyPair) error
	AddKeypair(parentSN string, keypair *KeyPair) error
	GetKeypair(SN string) (*KeyPair, error)
	GetFullPki() (*PkiNode, error)
}
