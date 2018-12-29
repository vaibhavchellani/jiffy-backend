package src

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"strings"
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
	Name    string  `json:"name"`
	Address string  `json:"address"`
	Network string  `json:"network"`
	ABI     abi.ABI `json:"abi"`
}

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

	abiBytes, err := MarshallABI(m.ABI)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	// convert address/name to lowercase to simplify search
	contract := ContractObj{
		Name:      m.Name,
		Address:   strings.ToLower(m.Address),
		Network:   m.Network,
		ABI:       abiBytes,
		QueryName: strings.ToLower(m.Name),
	}

	ControllerLogger.Debug("Contract registration initiated", "Address", contract.Address, "Name", contract.Name, "Network", contract.Network)

	err = c.DB.RegisterContract(contract)
	if err != nil {
		ControllerLogger.Error("Unable to register contract", "Error", err)
	}
	result, err := json.Marshal(map[string]interface{}{"status": "Success"})
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

func (c *Controller) GetContracts(w http.ResponseWriter, r *http.Request) {
	contracts, err := c.DB.GetContracts()
	if err != nil {
		ControllerLogger.Error("Unable to get all contracts", "Error", err)
	}
	result, err := json.Marshal(&contracts)
	if err != nil {
		ControllerLogger.Error("Error while marshalling get contracts response", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	ControllerLogger.Info("Successfully fetched all contracts", "Result", result)
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

// get a contract by name
func (c *Controller) GetContract(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//name := vars["name"]
	// filter contract by name/address
	filter := r.FormValue("filter")
	var contract ContractObj
	var err error
	// if filter is an address get contract by address else by name
	if common.IsHexAddress(filter) {
		contract, err = c.DB.GetContractViaAddr(filter)
		if err != nil {
			// TODO match err and sent respective status
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		contract, err = c.DB.GetContractViaName(filter)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	result, err := json.Marshal(&contract)
	if err != nil {
		ControllerLogger.Error("Error while marshalling get contract response", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	ControllerLogger.Info("Successfully fetched contract", "Result", result)
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}
