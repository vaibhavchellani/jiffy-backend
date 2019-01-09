package mongo

import (
	"github.com/globalsign/mgo"
	"github.com/jiffy-backend/config"
	"github.com/jiffy-backend/helper"
	"github.com/globalsign/mgo/bson"
)
type ILabelService interface {
	Register(label Label) error
	GetLabelByCreator(creator string,labels []Label) (err error)
}

// labels will refer functions by 4 byte unique signatures

type LabelService struct {
	collection *mgo.Collection
}

func labelModelIndex() mgo.Index {
	return mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
}

func NewLabelService(session *Session, dbName string) *LabelService {
	collection := session.GetCollection(dbName, config.LabelCollection)
	collection.EnsureIndex(labelModelIndex())
	return &LabelService{collection}
}

func (c *LabelService) Register(label Label) error {
	err := c.collection.Insert(label)
	if err != nil {
		return err
	}
	return nil
}


func (c *LabelService) GetLabelByCreator(creator string,labels *[]Label) (err error){
	err =c.collection.Find(bson.M{"creator_addr":creator}).All(labels)
	if err!=nil{
		helper.DBLogger.Error("Error fetching all labels","Creator",creator)
		return err
	}
	return nil
}


func (c *LabelService) GetLabelByID(_id bson.ObjectId,label Label) (err error) {
	err =c.collection.Find(bson.M{"_id":_id}).One(label)
	if err != nil {
		return err
	}
	return nil
}

func (c *LabelService) GetLabelByIDS(_ids []bson.ObjectId) ([]Label ,error) {
	var labels []Label
	for _,_id := range _ids {
		var label Label
		err :=c.collection.Find(bson.M{"_id":_id}).One(label)
		if err != nil {
			return labels,err
		}
		labels = append(labels, label)
	}
	return labels,nil
}


func (c *LabelService) GetLabelByContractName(contract string,labels *[]Label) (err error) {
	err =c.collection.Find(bson.M{"contract_name":contract}).All(labels)
	if err!=nil{
		helper.DBLogger.Error("Error fetching all labels","Creator",contract)
		return err
	}
	return nil
}

// add new function to existing label/modify
func (c *LabelService) AddFunctionToLabel(labelID bson.ObjectId,functions []Function) (err error) {
	selector:= bson.M{"_id":labelID}

	changeLog,err:=c.collection.Upsert(selector,bson.M{"$push":bson.M{"functions":functions}})
	if err!=nil{
		return err
	}
	helper.DBLogger.Debug("Added new functions to contract","ChangeInfo",changeLog,"Functions",functions)
	return nil
}


// will provide list of label ids
// will register new label with functions of prev labels in order
// will provide index where we need to add new functions if any
func (c *LabelService) MergeLabelsAndRegister(labelIDs []bson.ObjectId,newLabel Label,index int) (err error) {
	var functions []Function
	labels,err:=c.GetLabelByIDS(labelIDs)
	if err!=nil{
		return err
	}
	helper.DBLogger.Debug("Obtained all labels","LabelIDS",labelIDs,"Labels",labels)
	for i,label := range labels{
		if i !=index {
			helper.DBLogger.Debug("Inserting old label functions","Index",i,"Functions",label.Functions,"Count",len(label.Functions))
			functions = append(functions,label.Functions...)
		}else{
			helper.DBLogger.Debug("Inserting label functions in the index","Index",index,"Functions",newLabel.Functions,"Count",len(newLabel.Functions))
			functions = append(functions,newLabel.Functions...)
		}
	}
	helper.DBLogger.Debug("All functions arranged","Functions",functions,"Count",len(functions))
	newLabel.Functions = functions
	err =c.Register(newLabel)
	if err!=nil{
		return err
	}
	helper.DBLogger.Debug("Inserted new label with merged functions","IDs",labelIDs,"Label",newLabel)
	return nil
}







