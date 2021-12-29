package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type IndexEntry struct {
	SerialNumber string        `json:"sn"`
	Path         string        `json:"f"`
	Children     []*IndexEntry `json:"c,omitempty"`
	CA           bool          `json:"ca"`
}

type FilesystemStorer struct {
	BasePath string
}

func (idx *IndexEntry) toPkiNode() *PkiNode {
	children := []*PkiNode{}

	for i := range idx.Children {
		children = append(children, idx.Children[i].toPkiNode())
	}

	return &PkiNode{
		SerialNumber: idx.SerialNumber,
		Ca:           idx.CA,
		Children:     children,
	}
}

func (idx *IndexEntry) findChildSN(SN string) *IndexEntry { // recursivity FTW
	if idx.SerialNumber == SN {
		return idx
	}
	for _, child := range idx.Children {
		if found := child.findChildSN(SN); found != nil {
			return found
		}
	}
	return nil
}

func (storer *FilesystemStorer) getIndex() (*IndexEntry, error) {
	indexBytes, err := ioutil.ReadFile(filepath.Join(storer.BasePath, "index.json"))
	if err != nil {
		return nil, fmt.Errorf("opening index.json failed, %s", err)
	}

	var index IndexEntry
	if err := json.Unmarshal(indexBytes, &index); err != nil {
		return nil, fmt.Errorf("parsing index.json failed, %s", err)
	}

	return &index, nil
}

func (storer *FilesystemStorer) writeIndex(index *IndexEntry) error {
	indexBytes, err := json.Marshal(index)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(storer.BasePath, "index.json"), indexBytes, 0600); err != nil {
		return err
	}

	return nil
}

func (storer *FilesystemStorer) InitPki(name string, keypair *KeyPair) error {
	// Check if the folder already exists
	if _, err := os.Stat(storer.BasePath); !os.IsNotExist(err) {
		return fmt.Errorf("this PKI already exists")
	}

	if err := os.MkdirAll(storer.BasePath, 0700); err != nil {
		return fmt.Errorf("creating %s failed, %s", storer.BasePath, err)
	}

	if err := ioutil.WriteFile(filepath.Join(storer.BasePath, fmt.Sprintf("%s.crt", keypair.CertSN)), []byte(keypair.CertPem), 0600); err != nil {
		return fmt.Errorf("writing %s.crt failed, %s", keypair.CertSN, err)
	}

	if err := ioutil.WriteFile(filepath.Join(storer.BasePath, fmt.Sprintf("%s.key", keypair.CertSN)), []byte(keypair.KeyPem), 0600); err != nil {
		return fmt.Errorf("writing %s.key failed, %s", keypair.CertSN, err)
	}

	if err := storer.writeIndex(&IndexEntry{
		SerialNumber: keypair.CertSN,
		Path:         ".",
		Children:     []*IndexEntry{},
		CA:           true,
	}); err != nil {
		return fmt.Errorf("could not write index file, %s", err)
	}

	return nil
}

func (storer *FilesystemStorer) AddSubCA(parentSN string, keypair *KeyPair) error {
	// Find parent path
	index, err := storer.getIndex()
	if err != nil {
		return fmt.Errorf("failed to get the index, %s", err)
	}

	parentIndex := index.findChildSN(parentSN)
	if parentIndex == nil {
		return fmt.Errorf("could not find parent SN %s", parentSN)
	}

	// Check if the folder already exists
	subPath := filepath.Join(storer.BasePath, parentIndex.Path, keypair.CertSN)
	if _, err := os.Stat(subPath); !os.IsNotExist(err) {
		return fmt.Errorf("this subCA already exists")
	}

	if err := os.MkdirAll(subPath, 0700); err != nil {
		return fmt.Errorf("creating %s failed, %s", subPath, err)
	}

	// Write the certificate & key
	if err := ioutil.WriteFile(filepath.Join(subPath, fmt.Sprintf("%s.crt", keypair.CertSN)), []byte(keypair.CertPem), 0600); err != nil {
		return fmt.Errorf("writing %s.crt failed, %s", keypair.CertSN, err)
	}

	if err := ioutil.WriteFile(filepath.Join(subPath, fmt.Sprintf("%s.key", keypair.CertSN)), []byte(keypair.KeyPem), 0600); err != nil {
		return fmt.Errorf("writing %s.key failed, %s", keypair.CertSN, err)
	}

	// Update the index
	parentIndex.Children = append(parentIndex.Children, &IndexEntry{
		SerialNumber: keypair.CertSN,
		Path:         filepath.Join(parentIndex.Path, keypair.CertSN),
		Children:     []*IndexEntry{},
		CA:           true,
	})

	if err := storer.writeIndex(index); err != nil {
		return fmt.Errorf("could not write index file, %s", err)
	}

	return nil
}

func (storer *FilesystemStorer) GetKeypair(SN string) (*KeyPair, error) {
	index, err := storer.getIndex()
	if err != nil {
		return nil, fmt.Errorf("failed to get the index, %s", err)
	}

	entry := index.findChildSN(SN)
	if entry == nil {
		return nil, fmt.Errorf("failed to get SN %s", SN)
	}

	certPem, err := ioutil.ReadFile(filepath.Join(storer.BasePath, entry.Path, fmt.Sprintf("%s.crt", SN)))
	if err != nil {
		return nil, fmt.Errorf("failed to read %s.crt", SN)
	}

	keyPem, err := ioutil.ReadFile(filepath.Join(storer.BasePath, entry.Path, fmt.Sprintf("%s.key", SN)))
	if err != nil {
		return nil, fmt.Errorf("failed to read %s.key", SN)
	}

	return &KeyPair{
		CertPem: fmt.Sprintf("%s", certPem),
		KeyPem:  fmt.Sprintf("%s", keyPem),
		CertSN:  SN,
	}, nil
}

func (storer *FilesystemStorer) GetFullPki() (*PkiNode, error) {
	index, err := storer.getIndex()
	if err != nil {
		return nil, fmt.Errorf("failed to get the index, %s", err)
	}

	return index.toPkiNode(), nil
}
