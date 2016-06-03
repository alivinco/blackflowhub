package db
import "gopkg.in/mgo.v2"
import (
	"gopkg.in/mgo.v2/bson"
	"github.com/alivinco/blackflowhub/model"
	"fmt"
	"errors"
)

type MongodbAppMetaStore struct {
	session *mgo.Session
	db *mgo.Database
	appC *mgo.Collection
}

func (mg *MongodbAppMetaStore) Connect(serverName string , dbName string) error {
	session ,err := mgo.Dial(serverName)
	mg.session = session
	if err == nil {
		mg.session.SetMode(mgo.Monotonic, true)
		mg.db = mg.session.DB(dbName)
		mg.appC = mg.db.C("app_meta")
	}
	return err
}

func (mg *MongodbAppMetaStore) GetDbConnection()(*mgo.Session,*mgo.Database){
	return mg.session,mg.db
}
func (mg *MongodbAppMetaStore) Close (){
	mg.session.Close()
}

func (mg *MongodbAppMetaStore) Upsert(metaApp *model.App) (string,error){
	//TODO:use id if specified
	selector := bson.M{"appname":metaApp.AppName,"version":metaApp.Version,"developer":metaApp.Developer}
	info , err := mg.appC.Upsert(selector,*metaApp)
	if info.UpsertedId != nil {
		return info.UpsertedId.(bson.ObjectId).Hex(), err
	}else {
		return "",err
	}
}

func (mg *MongodbAppMetaStore) UpdateById(ID string ,metaApp *model.App) error{
	oid , err := mg.IdToObjectId(ID)
	if err == nil {
		_, err = mg.appC.UpsertId(oid, *metaApp)
	}
	return err
}

func (mg *MongodbAppMetaStore) Delete(ID string) error{
	return mg.appC.RemoveId(bson.ObjectIdHex(ID))
}

func (mg *MongodbAppMetaStore) IdToObjectId(ID string) (oid bson.ObjectId,err error) {

	defer func(){
		if r := recover(); r != nil{
			fmt.Println(r)
			fmt.Println("Invalid Id")
			err = errors.New("Invalid Id")
		}else {
			fmt.Println("Id is OK .")
		}

	}()

	return bson.ObjectIdHex(ID) , err
}

func (mg *MongodbAppMetaStore) GetById(ID string) (appResults *model.App, err error){
    appResults = &model.App{}
	oid , err := mg.IdToObjectId(ID)
	if err == nil {
		err = mg.appC.FindId(oid).One(appResults)
		fmt.Println("GetById results : ", appResults)
	}
	return appResults,err
}

func (mg *MongodbAppMetaStore) GetByFullName(developer string,appName string, version float32) (appResults *model.App, err error){
    appResults = &model.App{}
	selector := bson.M{"appname":appName,"version":version,"developer":developer}
	if err == nil {
		err = mg.appC.Find(selector).One(appResults)
		fmt.Println("GetByFullName results : ", appResults)
	}
	return appResults,err
}
// GetList returns list of all apps
func (mg*MongodbAppMetaStore) GetAll(filter *model.AppFilter) (*[]model.App,error){
	var appResults []model.App
	selector := bson.M{}
	if filter.Developer != "" {
	  selector = bson.M{"$or":[]bson.M{bson.M{"developer":filter.Developer},bson.M{"public":true}}}
	}else {
	  selector = bson.M{"public":filter.Public}
	}
 //	pipeline := bson.M{
//    "key1": 1,
//    "$or": []interface{}{
//        bson.M{"key2": 2},
//        bson.M{"key3": 2},
//    },
//  }

	err := mg.appC.Find(selector).All(&appResults)
	fmt.Println("Results All: ", appResults)
	return &appResults,err
}
