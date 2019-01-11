package helper

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/jiffy-backend/config"
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

func GenerateHash(network string, address string) [32]byte {
	identifier := fmt.Sprintf(strings.ToLower(network), strings.ToLower(address))
	return sha256.Sum256([]byte(identifier))
}

func GetNetworkDetails(URL string) (name string, err error) {
	switch URL {
	case config.MainNetChainURL:
		return "main", err
	case config.RopstenChainURL:
		return "ropsten", err
	case config.RinkelbyChainURL:
		return "rinkelby", err
	case config.KovanChainURL:
		return "kovan", err
	default:
		// TODO Check if URL is valid URL
		// generate err if not valid URL
		return "custom", err
	}
}

// get JSON error got request
func Error(w http.ResponseWriter, code int, message string) {
	JSONResponse(w, code, map[string]string{"error": message})
}

// get JSON response for request
func JSONResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// validates valid address
func IsValidAddress(address string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

	if re.MatchString(address) {
		return false
	}
	return true
}

// get tag of struct at an index
func GetModelFieldAtIndex(v interface{}, index int) string {
	st := reflect.TypeOf(v)
	field := st.Field(index)
	return field.Tag.Get("bson")

}
