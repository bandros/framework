package framework

import (
	"context"
)

type DatabaseMongoDb struct {
	Ctx   context.Context
	table string
}

func (p *DatabaseMongoDb) Table(tbl string) *DatabaseMongoDb {
	p.table = tbl
	return p
}

func (p *DatabaseMongoDb) Insert(document interface{}) (interface{}, error) {
	var db, err = MongoDbConnect(p.Ctx)
	if err != nil {
		return nil, err
	}
	i, err := db.Collection(p.table).InsertOne(p.Ctx, document)
	if err != nil {
		return nil, err
	}
	return i.InsertedID, nil
}

func (p *DatabaseMongoDb) InsertBatch(documents []interface{}) ([]interface{}, error) {
	var db, err = MongoDbConnect(p.Ctx)
	if err != nil {
		return nil, err
	}
	i, err := db.Collection(p.table).InsertMany(p.Ctx, documents)
	if err != nil {
		return nil, err
	}
	return i.InsertedIDs, nil
}
