package eth

import (
	"Share-Wallet/pkg/wallet/account"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

// Note: key filename is different between Geth and Parity
// Geth
// - 0x71678cd07cfac46c2dc427f999abf46aae115925
// - UTC--2018-10-11T06-58-43.857846090Z--71678cd07cfac46c2dc427f999abf46aae115925

// Parity, filename includes just UUID
// - "0xcf9583c3c10cf895af95a2810243765c4fe7c038",
// - UTC--2018-10-11T06-59-28Z--2bd02735-84ec-593e-f2b2-73cce1b1862c

// File names for Parity keys
// https://ethereum.stackexchange.com/questions/13951/file-names-of-parity-keys

// So Parity key filename should be renamed to same format as Geth
// e.g. UTC--2018-10-12T01-53-58Z--fff7e98d-b3b7-08f4-65cd-3fe82416cebf--45783b86c2aa1ce81632ac2db26a91acc3ea6226

// ToECDSA converts privKey to ECDSA
func (e *Ethereum) ToECDSA(privKey string) (*ecdsa.PrivateKey, error) {
	var bytePrivKey []byte
	var err error
	if strings.HasPrefix(privKey, "0x") || strings.HasPrefix(privKey, "0X") {
		bytePrivKey, err = hexutil.Decode(privKey)
	} else {
		bytePrivKey, err = hex.DecodeString(privKey)
	}
	if err != nil {
		return nil, errors.Wrap(err, "fail to call hex.DecodeString()")
	}
	return crypto.ToECDSA(bytePrivKey)
}

// GetKeyDir returns keystore directory
func (e *Ethereum) GetKeyDir() string {
	return e.keyDir
}

// GetPrivKey returns keystore.Key object
func (e *Ethereum) GetPrivKey(hexAddr, password string) (*keystore.Key, error) {
	keyDir := e.GetKeyDir()
	keyJSON, err := e.readPrivKey(hexAddr, keyDir)
	if err != nil {
		log.Printf("err:%s", err.Error())
		return nil, errors.Wrapf(err, "fail to call e.readPrivKey()")
	}
	if keyJSON == nil {
		// file is not found
		log.Println("file is not found")
		return nil, errors.New("private key file is not found")
	}

	key, err := keystore.DecryptKey(keyJSON, password)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to call keystore.DecryptKey()")
	}
	log.Print(key.PrivateKey)

	return key, nil
}

// readPrivKey read private key file from directory
// Note: file is found out from local directory,
//  if node is working remotely, file is not found.
func (e *Ethereum) readPrivKey(hexAddr, path string) ([]byte, error) {
	// search file
	// filename is like `UTC--2020-05-18T16-01-32.772616000Z--e52307deb1a7dc3985d2873b45ae23b91d57a36d`
	//                   UTC--2020-05-22T13-59-51.645058000Z--d4ec46122b3f0afc0287144adcca5d65b22b799c
	// Note: all letter of address in filename is a lowercase letter
	addr := strings.TrimPrefix(strings.ToLower(hexAddr), "0x")

	files, err := filepath.Glob(fmt.Sprintf("%s/*--%s", path, addr))
	if err != nil {
		return nil, errors.Wrap(err, "fail to call filepath.Glob()")
	}
	if len(files) == 0 {
		// file is not found
		return nil, errors.New("private key file is not found")
	}
	if len(files) > 1 {
		return nil, errors.Errorf("target private key files are found more than 1 by %s", addr)
	}
	log.Printf("file path:%s", files[0])

	return ioutil.ReadFile(files[0])
}

// RenameParityKeyFile renames parity file format
// Deprecated
func (e *Ethereum) RenameParityKeyFile(hexAddr string, accountType account.AccountType) error {
	if !e.isParity {
		return nil
	}

	files, err := ioutil.ReadDir(e.GetKeyDir())
	if err != nil {
		return err
	}

	var fileNames []string
	for _, v := range files {
		if v.IsDir() {
			continue
		}
		fileNames = append(fileNames, v.Name())
	}
	sort.Strings(fileNames)

	// get last one
	target := fileNames[len(fileNames)-1]

	// remove `0x`を from hexAddr
	addr := strings.TrimLeft(hexAddr, "0x")

	// rename xxxxx--[address]
	previousName := fmt.Sprintf("%s/%s", e.GetKeyDir(), target)
	os.Rename(previousName, fmt.Sprintf("%s--%s", previousName, addr))

	return nil
}
