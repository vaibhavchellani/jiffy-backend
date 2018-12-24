package src

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"encoding/json"
)

func MarshallABI(abi abi.ABI) ([]byte,error){
	abiBytes, err := json.Marshal(abi)
	if err != nil {
		return abiBytes,err
	}
	return abiBytes,nil
}

func UnMarshallABI(abiBytes []byte,abi abi.ABI) error {
	err:=json.Unmarshal(abiBytes,&abi)
	if err!=nil{
		return err
	}
	return nil
}