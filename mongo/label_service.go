package mongo

import "github.com/globalsign/mgo"

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

func NewLabelService(session *Session, dbName string, collectionName string) *LabelService {
	collection := session.GetCollection(dbName, collectionName)
	collection.EnsureIndex(contractModelIndex())
	return &LabelService{collection}
}
