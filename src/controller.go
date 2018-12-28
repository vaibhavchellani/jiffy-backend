package src

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
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

	contract := ContractObj{
		Name:    m.Name,
		Address: m.Address,
		Network: m.Network,
		ABI:     abiBytes,
	}
	ControllerLogger.Debug("Contract registration initiated", "Address", contract.Address, "Name", contract.Name, "Network", contract.Network)

	err = c.DB.RegisterContract(contract)
	if err != nil {
		ControllerLogger.Error("Unable to register contract", "Error", err)
	}
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
	w.Write(result)
}

func (c *Controller) GetContract(w http.ResponseWriter, r *http.Request) {
	conrtact, err := c.DB.GetContract()
	result, err := json.Marshal(&conrtact)
	if err != nil {
		ControllerLogger.Error("Error while marshalling get contract response", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	ControllerLogger.Info("Successfully fetched contract", "Result", result)
	w.Write(result)

}
