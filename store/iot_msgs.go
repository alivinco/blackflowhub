package store

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/alivinco/blackflowhub/model"
	"fmt"
	"time"
)

type IotMsgStore struct {
	session *mgo.Session
	db *mgo.Database
	iotMsgC *mgo.Collection
}

func NewIotMsgStore(session *mgo.Session,db *mgo.Database)(*IotMsgStore){
	imst := IotMsgStore{session:session,db:db}
	imst.iotMsgC = db.C("iot_msg")
	return &imst
}

func (ms *IotMsgStore) UpsertMsg(iotMsg *model.IotMsg) (string,error){
	var selector bson.M
	if len(iotMsg.ID)>0 {
		selector = bson.M{"_id":iotMsg.ID}
		iotMsg.Created = time.Now()
		iotMsg.Updated = time.Now()
	}else{
		selector = bson.M{"payloadtype":iotMsg.PayloadType,"type":iotMsg.Type,"class":iotMsg.Class,"subclass":iotMsg.SubClass,"version":iotMsg.Version}
		iotMsg.Updated = time.Now()
	}
	info , err := ms.iotMsgC.Upsert(selector,*iotMsg)
	if err == nil {
		if info.UpsertedId != nil {
			return info.UpsertedId.(bson.ObjectId).Hex(), err
		} else {
			return "", err
		}
	} else {
		return "" , err
	}

}

func (ms *IotMsgStore) DeleteMsg(ID string) error{
	return ms.iotMsgC.RemoveId(bson.ObjectIdHex(ID))
}

// GetList returns list of all apps
func (ms *IotMsgStore) GetIotMsgs(filter *model.IotMsg) (*[]model.IotMsg,error){
	var results []model.IotMsg
	selector := bson.M{}
	if len(filter.ID)>0 {
		selector = bson.M{"_id":filter.ID}
	}else if filter.Author != "" {
		selector = bson.M{"$or":[]bson.M{bson.M{"author":filter.Author},bson.M{"public":true}}}
	}else {
		selector = bson.M{"public":filter.Public}
	}
	err := ms.iotMsgC.Find(selector).All(&results)
	fmt.Println("Results All: ", results)
	return &results,err
}
