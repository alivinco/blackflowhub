package main
import "github.com/gin-gonic/gin"
import "flag"
import (
	"github.com/alivinco/blackflowhub/controller"
	"github.com/alivinco/blackflowhub/store"
	"github.com/alivinco/blackflowhub/store/db"
	"github.com/alivinco/blackflowhub/gincontrib/jwt"
//	"github.com/auth0/go-jwt-middleware"
//	"github.com/dgrijalva/jwt-go"
	"fmt"
	"encoding/base64"
	"github.com/itsjamie/gin-cors"
	"time"
	"os"
)

var mongoDb db.MongodbAppMetaStore
var fileStore db.FileStore
var appStore store.AppStore

func initAppStore(dbConn string , dbName string , fsLocation string) (error){
	err := mongoDb.Connect(dbConn,dbName)
	if err==nil {
		fileStore = db.FileStore{fsLocation}
		appStore = store.AppStore{&mongoDb, &fileStore}
		return nil
	}else {
		return err
	}

}

func stopAppStore(){
	mongoDb.Close()
}



func RunHttpServer(bindAddress string,jwtSecret string){
	r := gin.Default()
	r.Use(cors.Middleware(cors.Config{
				Origins:        "*",
				Methods:        "GET, PUT, POST, DELETE",
				RequestHeaders: "Origin, Authorization, Content-Type",
				ExposedHeaders: "",
				MaxAge: 50 * time.Second,
				Credentials: true,
				ValidateHeaders: false,
			}))
	decoded_secret, _ := base64.URLEncoding.DecodeString(jwtSecret)
	r.Use(jwt.Auth(string(decoded_secret)))

	contr := controller.AppRestController{&appStore}
	root := r.Group("/bfhub/api")
	{
		// public and logged in users
		root.GET("/apps",contr.GetApps)
		root.GET("/apps/id/:app_id",contr.GetAppById)
		root.GET("/apps/by_full_name",contr.GetAppByFullName)
		root.GET("/apps/file/:app_id",contr.GetFile)
		root.POST("/apps",contr.PostApp)
		root.DELETE("/apps/:app_id",contr.DeleteApps)
		//curl -i -v -F "ID=test_file" -F "file=@package.json" http://localhost:8080/bfhub/api/file
		root.POST("/apps/file",contr.PostFile)
	}
	r.Run(bindAddress)
}

func main(){
	bindAddress := ":5050"
	dbConn := "localhost"
	dbName := "blackflow"
	fsLocation := "data/"
	jwtSecret := ""
    // Load configs from env variable or from command line .
	bindAddress = os.Getenv("BFH_BIND_ADDR")
	if bindAddress != "" {
		dbConn = os.Getenv("BFH_DB_CONN_STR")
		dbName = os.Getenv("BFH_DB_NAME")
		fsLocation = os.Getenv("BFH_FS_LOCATION")
		jwtSecret = os.Getenv("BFH_JWT_SECRET")
	}else{
		flag.StringVar(&bindAddress,"addr",":5050","Server bind address")
		flag.StringVar(&dbConn,"db_conn","localhost","Mongo db connection string.Default = localhost. Example mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb")
		flag.StringVar(&dbName,"db_name","blackflow","Database name")
		flag.StringVar(&fsLocation,"fs_location","data/","File store location")
		flag.StringVar(&jwtSecret,"jwt_secret","","Jwt secret")
	}

	flag.Parse()
	fmt.Println("addr:",bindAddress)
	fmt.Println("db_conn:",dbConn)
	fmt.Println("db_name:",dbName)
	fmt.Println("fs_location:",fsLocation)
	fmt.Println("jwt_secret:",jwtSecret)

	err := initAppStore(dbConn,dbName,fsLocation)
	if err == nil {
		RunHttpServer(bindAddress, jwtSecret)
		defer stopAppStore()
	}else{
		fmt.Println("App can't be started because erros",err)
	}

}