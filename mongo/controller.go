package mongo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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
	ABI     string `json:"abi"`               // abi of contract
	Owner   string  `json:"owner"`             // owner aka address registering the contract
	Network string  `json:"network"` // network URL
}

// TODO make sure name doesnt match with existing routes
// handler for contract registration
func (c *Controller) RegisterContract(w http.ResponseWriter, r *http.Request) {
	var m ContractInput
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helper.Error(w,http.StatusBadRequest,err.Error())
		return
	}
	err = json.Unmarshal(body, &m)
	if err != nil {
		helper.Error(w,http.StatusBadRequest,err.Error())
		return
	}


	name, err := helper.GetNetworkDetails(m.Network)
	if err != nil {
		helper.Error(w,http.StatusBadRequest,err.Error())
		return
	}

	contractHash:=helper.GenerateHash(name,m.Address)
	contractHashStr:=hex.EncodeToString(contractHash[:])

	// convert address/name to lowercase to simplify search
	contract := ContractObj{
		Name:        m.Name,
		Address:     strings.ToLower(m.Address),
		NetworkName: name,
		ABI:         m.ABI,
		QueryName:   strings.ToLower(m.Name),
		Owner:       strings.ToLower(m.Owner),
		Identifier:  contractHashStr,
		NetworkURL:  m.Network,
	}

	helper.ControllerLogger.Debug("hash generated","hash",contractHashStr)
	// link to exisitng contract by name
	existingContract,err:=c.DB.GetContractByIdentifier(contractHashStr)
	if err == nil{
		if strings.Compare(existingContract.Cloned,"")==0{
			contract.Cloned = existingContract.Name
		}else{
			contract.Cloned = existingContract.Cloned
		}

	}

	if err:=contract.ValidateBasic(); err!=nil {
		// add error
		return
	}

	helper.ControllerLogger.Debug("Contract registration initiated", "Address", contract.Address, "Name", contract.Name, "Network", contract.NetworkName)

	err = c.DB.RegisterContract(contract)
	if err != nil {
		helper.ControllerLogger.Error("Unable to register contract", "Error", err)
		helper.Error(w,http.StatusBadRequest,err.Error())
		return
	}

	helper.JsonResponse(w, http.StatusOK, map[string]interface{}{"status": "Success", "Contract": contract.Json()})
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
		helper.Error(w,http.StatusBadRequest,err.Error())
	}
	helper.JsonResponse(w, http.StatusOK, map[string]interface{}{"status": "Success", "Contract": contracts})
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
			helper.Error(w,http.StatusBadRequest,err.Error())
			return
		}
	} else {
		contract, err = c.DB.GetContractByName(filter)
		if err != nil {
			helper.Error(w,http.StatusBadRequest,err.Error())
			return
		}
	}

	helper.JsonResponse(w, http.StatusOK, map[string]interface{}{"status": "Success", "Contract": contract})
}

// get dapp given dapp name and spit abi
func (c *Controller) GetDapp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["dapp_name"]
	contract, err := c.DB.GetContractByName(name)
	if err != nil {
		helper.Error(w,http.StatusBadRequest,err.Error())
		return
	}
	helper.JsonResponse(w, http.StatusOK, map[string]interface{}{"status": "Success", "Contract":contract.ABI})
}

func (c *Controller) CheckExistence(w http.ResponseWriter, r *http.Request) {
	addr := r.FormValue("address")
	network := r.FormValue("network")
	hash := helper.GenerateHash(network,addr)
	hashStr:= hex.EncodeToString(hash[:])
	contract, err := c.DB.GetContractByIdentifier(hashStr)
	if err != nil {
		helper.Error(w,http.StatusBadRequest,err.Error())
		return
	}
	helper.JsonResponse(w, http.StatusOK, map[string]interface{}{"status": "Success", "Contract": contract.Json()})
}

