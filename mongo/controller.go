package mongo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/jiffy-backend/helper"
	"github.com/pkg/errors"
	"strings"
	"encoding/hex"
)

type Controller struct {
	DB DB
}

// handler for registrations (labels / contracts)
func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	entity := vars["entity"]
	if entity == "contract" {
		c.RegisterContract(w, r)
	} else if entity == "label" {
		c.RegisterLabel(w, r)
	} else {
		err := errors.New("Invalid Entity Registration")
		w.Write([]byte(err.Error()))
		return
	}
}

type ContractInput struct {
	Name    string  `json:"name"`              // name given to contract => used to generate unique URL
	Address string  `json:"address"`           // address of contract
	ABI     abi.ABI `json:"abi"`               // abi of contract
	Owner   string  `json:"owner"`             // owner aka address registering the contract
	Network string  `json:"network"` // network URL
}

// TODO make sure name doesnt match with existing routes
// handler for contract registration
func (c *Controller) RegisterContract(w http.ResponseWriter, r *http.Request) {
	var m ContractInput
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = json.Unmarshal(body, &m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	abiBytes, err := helper.MarshallABI(m.ABI)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// WIP
	//var newAbi abi.ABI
	////helper.UnMarshallABI(abiBytes,&newAbi)
	//if err := newAbi.UnmarshalJSON(abiBytes); err != nil {
	//	helper.ControllerLogger.Error("unable to unmarshall")
	//}
	//// TODO figure out a way to save abi and reproduce the same
	//helper.ControllerLogger.Debug("new abi", "bytes", abiBytes, "unmae", newAbi.UnmarshalJSON(abiBytes), "Inp", m.ABI)


	name, err := helper.GetNetworkDetails(m.Network)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// convert address/name to lowercase to simplify search
	contract := ContractObj{
		Name:        m.Name,
		Address:     strings.ToLower(m.Address),
		NetworkName: name,
		ABI:         abiBytes,
		QueryName:   strings.ToLower(m.Name),
		Owner:       strings.ToLower(m.Owner),
		Identifier:  helper.GenerateHash(name,m.Address),
		NetworkURL:  m.Network,
	}

	helper.ControllerLogger.Debug("hash generated","hash",hex.EncodeToString(contract.Identifier[:]))
	if err:=contract.ValidateBasic(); err!=nil {
		// add error
		return
	}

	helper.ControllerLogger.Debug("Contract registration initiated", "Address", contract.Address, "Name", contract.Name, "Network", contract.NetworkName)

	err = c.DB.RegisterContract(contract)
	if err != nil {
		helper.ControllerLogger.Error("Unable to register contract", "Error", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	result, err := json.Marshal(map[string]interface{}{"status": "Success", "Contract": contract.Json()})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(result))
}

// handler for label registration
func (c *Controller) RegisterLabel(w http.ResponseWriter, r *http.Request) {
	//c.DB.RegisterContract()
}

// get all contracts
func (c *Controller) GetContracts(w http.ResponseWriter, r *http.Request) {
	contracts, err := c.DB.GetContracts()
	if err != nil {
		helper.ControllerLogger.Error("Unable to get all contracts", "Error", err)
	}
	result, err := json.Marshal(&contracts)
	if err != nil {
		helper.ControllerLogger.Error("Error while marshalling get contracts response", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	helper.ControllerLogger.Info("Successfully fetched all contracts", "Result", result)
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

// get a contract by name
func (c *Controller) GetContract(w http.ResponseWriter, r *http.Request) {
	// filter contract by name/address
	filter := r.FormValue("filter")
	var contract ContractObj
	var err error
	// if filter is an address get contract by address else by name
	if common.IsHexAddress(filter) {
		contract, err = c.DB.GetContractByAddr(filter)
		if err != nil {
			// TODO match err and sent respective status
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		contract, err = c.DB.GetContractByName(filter)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	result, err := json.Marshal(&contract)
	if err != nil {
		helper.ControllerLogger.Error("Error while marshalling get contract response", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

// get dapp given dapp name and spit abi
func (c *Controller) GetDapp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["dapp_name"]
	contract, err := c.DB.GetContractByName(name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(contract.ABI)
}

func (c *Controller) CheckExistence(w http.ResponseWriter, r *http.Request) {
	addr := r.FormValue("address")
	network := r.FormValue("network")
	hash := helper.GenerateHash(network,addr)

	contract, err := c.DB.GetContractByIdentifier(hash)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	helper.JsonResponse(w, http.StatusOK, map[string]interface{}{"status": "Success", "Contract": contract.Json()})
}

