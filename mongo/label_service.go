package mongo

import (
	"github.com/globalsign/mgo"
	"github.com/jiffy-backend/config"
)

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
	collection.EnsureIndex(contractModelIndex())
	return &LabelService{collection}
}

func (c *LabelService) Register(label Label) error {
	err := c.collection.Insert(label)
	if err != nil {
		return err
	}
	return nil
}
