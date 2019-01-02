package helper

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"fmt"
	"crypto/sha256"
)

func MarshallABI(abi abi.ABI) ([]byte, error) {
	abiBytes, err := json.Marshal(abi)
	if err != nil {
		return abiBytes, err
	}
	return abiBytes, nil
}

func UnMarshallABI(abiBytes []byte, abi *abi.ABI) error {
	err := json.Unmarshal(abiBytes, abi)
	if err != nil {
		return err
	}
	return nil
}

func GenerateHash(network string,address string) ([32]byte)  {
	identifier :=fmt.Sprintf(network,address)
	return sha256.Sum256([]byte(identifier))
}