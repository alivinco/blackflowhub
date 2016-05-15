package db
import (
	"path/filepath"
	"io/ioutil"
	"fmt"
	"os"
)

type FileStore struct {
	StoreDir string
}
func (fs *FileStore) SaveFile(appName string,binFile []byte) error  {
		 fullPath := filepath.Join(fs.StoreDir,appName)
		 fmt.Println("Writing file to",fullPath)
  		 return ioutil.WriteFile(fullPath,binFile,0777)
}

func (fs *FileStore) GetFile(appName string) ([]byte , error){
	fcontent , err := ioutil.ReadFile(filepath.Join(fs.StoreDir,appName))
	return fcontent,err
}

func (fs *FileStore) DeleteFile(appName string) (error){
	return os.Remove(filepath.Join(fs.StoreDir,appName))
}

func (fs *FileStore) GetStorageLocation() string{
	return fs.StoreDir
}