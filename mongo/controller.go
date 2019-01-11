package mongo

//
//
// All input validation to be performed here
//
//

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"encoding/hex"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/jiffy-backend/config"
	"github.com/jiffy-backend/helper"
	"github.com/pkg/errors"
)

// controller struct for DB access
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

// handler for updating (labels / contracts)
func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	entity := vars["entity"]
	if entity == "contract" {
		c.RegisterContract(w, r)
	} else if entity == "label" {
		c.RegisterLabel(w, r)
	} else {
		err := errors.New("Invalid Entity updation")
		w.Write([]byte(err.Error()))
		return
	}
}

// input for registering contract
type ContractInput struct {
	Name    string `json:"name"`    // name given to contract => used to generate unique URL
	Address string `json:"address"` // address of contract
	ABI     string `json:"abi"`     // abi of contract
	Owner   string `json:"owner"`   // owner aka address registering the contract
	Network string `json:"network"` // network URL
}

// TODO make sure name doesnt match with existing routes
// handler for contract registration
func (c *Controller) RegisterContract(w http.ResponseWriter, r *http.Request) {
	var m ContractInput
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	err = json.Unmarshal(body, &m)
	if err != nil {
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	name, err := helper.GetNetworkDetails(m.Network)
	if err != nil {
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	contractHash := helper.GenerateHash(name, m.Address)
	contractHashStr := hex.EncodeToString(contractHash[:])

	// convert address/name to lowercase to simplify search
	contract := ContractObj{
		Name:         m.Name,
		Address:      strings.ToLower(m.Address),
		NetworkName:  name,
		ABI:          m.ABI,
		QueryName:    strings.ToLower(m.Name),
		Owner:        strings.ToLower(m.Owner),
		Identifier:   contractHashStr,
		NetworkURL:   m.Network,
		TimeCreated:  time.Now(),
		LastModified: time.Now(),
	}

	helper.ControllerLogger.Debug("hash generated", "hash", contractHashStr)

	if err := contract.ValidateBasic(); err != nil {
		// add error
		return
	}

	// link to exisitng contract by name
	existingContract, err := c.DB.GetContractByIdentifier(contractHashStr)
	if err == nil {
		if strings.Compare(existingContract.Cloned, "") == 0 {
			contract.Cloned = existingContract.Name
		} else {
			contract.Cloned = existingContract.Cloned
		}
	}

	helper.ControllerLogger.Debug("Contract registration initiated", "Address", contract.Address, "Name", contract.Name, "Network", contract.NetworkName)

	err = c.DB.RegisterContract(contract)
	if err != nil {
		helper.ControllerLogger.Error("Unable to register contract", "Error", err)
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.JSONResponse(w, http.StatusOK, map[string]interface{}{"status": "Success", "Contract": contract.Json()})
}

// label input
type LabelInput struct {
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	CreatorAddr  string      `json:"creator"`
	Functions    []FuncInput `json:"functions"`
	ContractName string      `json:"contract_name"`
	MergeLabels  []string    `json:"merge_labels"`
	Index        int         `json:"label_index"`
}

// input for label functions
type FuncInput struct {
	MethodSig   string `json:"method_sig"`
	Skippable   bool   `json:"skippable"`
	Usage       int    `json:"usage"` // transaction or call(1 == tx , 0== call)
	Description string `json:"description"`
}

// handler for label registration
func (c *Controller) RegisterLabel(w http.ResponseWriter, r *http.Request) {
	var m LabelInput
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	err = json.Unmarshal(body, &m)
	if err != nil {
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	if strings.Compare(m.ContractName, "") == 0 {
		err = errors.New("Empty contract name not allowed")
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	contract, err := c.DB.GetContractByName(m.ContractName)
	if err != nil {
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	var functions []Function
	for _, functionInput := range m.Functions {
		usage := config.Call
		if functionInput.Usage == 1 {
			usage = config.Transaction
		}
		function := Function{
			MethodSig:   []byte(functionInput.MethodSig),
			Skippable:   functionInput.Skippable,
			Usage:       usage,
			Description: functionInput.Description,
		}
		functions = append(functions, function)
	}

	label := Label{
		ContractName: m.ContractName,
		ContractID:   contract.ID,
		CreatorAddr:  contract.Address,
		Functions:    functions,
		Name:         m.Name,
		Description:  m.Description,
		TimeCreated:  time.Now(),
	}
	err = c.DB.RegisterLabel(label)
	if err != nil {
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	helper.JSONResponse(w, http.StatusOK, map[string]interface{}{"status": "Success", "Label": label})
}

// get all contracts
func (c *Controller) GetContracts(w http.ResponseWriter, r *http.Request) {
	contracts, err := c.DB.GetContracts()
	if err != nil {
		helper.ControllerLogger.Error("Unable to get all contracts", "Error", err)
		helper.Error(w, http.StatusBadRequest, err.Error())
	}
	helper.JSONResponse(w, http.StatusOK, map[string]interface{}{"status": "Success", "Contract": contracts})
}

// get a contract by name
func (c *Controller) GetContract(w http.ResponseWriter, r *http.Request) {
	// filter contract by name/address
	filter := r.FormValue("filter")
	if strings.Compare(filter, "") == 0 {
		err := errors.New("Please provide address or name of contract to search")
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	var contract ContractObj
	var err error
	// if filter is an address get contract by address else by name
	if common.IsHexAddress(filter) {
		contract, err = c.DB.GetContractByAddr(filter)
		if err != nil {
			// TODO match err and sent respective status
			helper.Error(w, http.StatusBadRequest, err.Error())
			return
		}
	} else {
		contract, err = c.DB.GetContractByName(filter)
		if err != nil {
			helper.Error(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	helper.JSONResponse(w, http.StatusOK, map[string]interface{}{"status": "Success", "Contract": contract})
}

// get dapp given dapp name and spit abi
func (c *Controller) GetDapp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["dapp_name"]
	if strings.Compare(name, "") == 0 {
		// TODO should be redirected to home page
	}
	contract, err := c.DB.GetContractByName(name)
	if err != nil {
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	helper.JSONResponse(w, http.StatusOK, map[string]interface{}{"status": "Success", "Contract": contract.ABI})
}

// checks existence of contract by creating identifier = contract address + network
func (c *Controller) CheckExistence(w http.ResponseWriter, r *http.Request) {
	addr := r.FormValue("address")
	network := r.FormValue("network")
	if strings.Compare(addr, "") == 0 {
		err := errors.New("Please provide address")
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	if strings.Compare(network, "") == 0 {
		err := errors.New("Please provide network")
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	hash := helper.GenerateHash(network, addr)
	hashStr := hex.EncodeToString(hash[:])
	contract, err := c.DB.GetContractByIdentifier(hashStr)
	if err != nil {
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	helper.JSONResponse(w, http.StatusOK, map[string]interface{}{"status": "Success", "Contract": contract.Json()})
}

// -------- Label related controllers

// Gets label by contract name or address
func (c *Controller) GetLabelsByContract(w http.ResponseWriter, r *http.Request) {
	contract := r.FormValue("contract")
	if strings.Compare(contract, "") == 0 {
		err := errors.New("Please provide address or name of contract to search")
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	var labels []Label
	var err error
	if common.IsHexAddress(contract) {
		// get labels by contract address
		labels, err = c.DB.GetLabelViaContractAddr(contract)
		helper.Error(w, http.StatusBadRequest, err.Error())
		if err != nil {
			return
		}
	} else {
		// get labels by contract name
		labels, err = c.DB.GetLabelViaContractName(contract)
		if err != nil {
			helper.Error(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	helper.JSONResponse(w, http.StatusOK, map[string]interface{}{"status": "Success", "Labels": labels, "Count": len(labels)})
}

// get labels by creator address
func (c *Controller) GetLabelsByCreator(w http.ResponseWriter, r *http.Request) {
	creatorAddr := r.FormValue("creator")
	if strings.Compare(creatorAddr, "") == 0 {
		err := errors.New("Please provide address or name of contract to search")
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// TODO validate address to be propoper eth address

	labels, err := c.DB.GetLabelViaCreator(creatorAddr)
	if err != nil {
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.JSONResponse(w, http.StatusOK, map[string]interface{}{"status": "Success", "Labels": labels, "Count": len(labels)})
}

// --------
