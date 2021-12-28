package storage

import (
	"pkitool/models"
)

type KeyPair struct {
	CertPem string
	CertSN  string
	KeyPem  string
}

type Storer interface {
	InitPki(name string, keypair *KeyPair) error
	AddSubCA(parentSN string, keypair *KeyPair) error
	GetKeypair(SN string) (*KeyPair, error)
	GetFullPki() (*models.PkiDetails, error)
}
