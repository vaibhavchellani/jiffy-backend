package src

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type Controller struct {
	DB DB
}

var controllerLogger = Logger.With("module", "controller")

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

	abiBytes,err:=MarshallABI(m.ABI)
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
	controllerLogger.Debug("Contract registration initiated", "Address", contract.Address, "Name", contract.Name, "Network", contract.Network)

	err = c.DB.RegisterContract(contract)
	if err != nil {
		controllerLogger.Error("Unable to register contract", "Error", err)
	}
}

// handler for label registration
func (c *Controller) RegisterLabel(w http.ResponseWriter, r *http.Request) {
	//c.DB.RegisterContract()
}
