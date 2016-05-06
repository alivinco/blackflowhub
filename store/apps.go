package store
import (
	"github.com/alivinco/blackflowhub/model"
	"time"
)

type AppMetaStore interface {
	Upsert(*model.App) (string,error)
	UpdateById(string ,*model.App) error
	Delete(string) error
	GetAll(*model.AppFilter) (*[]model.App,error)
	GetById(string)(*model.App,error)
	GetByFullName(developer string,appName string,version float32)(*model.App,error)
//	Find(*model.AppFind)
}

type AppFileStore interface {
	SaveFile(string,[]byte) error
	GetFile(string) ([]byte , error)
	GetStorageLocation() string
	DeleteFile(string) error
}

type AppStore struct {
	AppMetaStore AppMetaStore
	AppFileStore AppFileStore
}


func(as *AppStore) UpsertApp(metaApp *model.App) (string,error){
	metaApp.Updated = time.Now()
	return as.AppMetaStore.Upsert(metaApp)

}

func(as *AppStore) GetApps(filter *model.AppFilter) (*[]model.App,error) {
	return as.AppMetaStore.GetAll(filter)
}

func(as *AppStore) GetAppById(ID string) (*model.App,error) {
	return as.AppMetaStore.GetById(ID)
}

func(as *AppStore) GetAppByFullName(developer string,appName string,version float32) (*model.App,error) {
	return as.AppMetaStore.GetByFullName(developer,appName,version)
}

func(as *AppStore) UpdateAppById(ID string,metaApp *model.App) error{
	return as.AppMetaStore.UpdateById(ID,metaApp)
}

func(as *AppStore) DeleteApp(ID string) error{
	err := as.AppMetaStore.Delete(ID)
	if err == nil {
		err = as.AppFileStore.DeleteFile(ID)
	}
	return err
}

func(as *AppStore) SaveFile(ID string,filename string, fileContent []byte) error {
	//TODO:Add ID validation against MetaStore and file name has to be added to metastore.
	am , err := as.AppMetaStore.GetById(ID)
	if err == nil {
		am.File = filename
		am.Updated = time.Now()
		am.Size = int32(len(fileContent))
		as.AppMetaStore.UpdateById(ID, am)
		return as.AppFileStore.SaveFile(ID, fileContent)
	}else {
		return err
	}

}
func(as *AppStore) GetFile(ID string, fileContent []byte) error {
	return as.AppFileStore.SaveFile(ID , fileContent)

}
