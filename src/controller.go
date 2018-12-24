package src

import "net/http"

type Controller struct {
	DB DB
}

func (c *Controller) RegisterContract(w http.ResponseWriter, r *http.Request){
	c.DB.GetAllContracts()
}