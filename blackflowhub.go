package main
import "github.com/gin-gonic/gin"
import "flag"
import (
	"github.com/alivinco/blackflowhub/controller"
	"github.com/alivinco/blackflowhub/store"
	"github.com/alivinco/blackflowhub/store/db"
	"github.com/alivinco/blackflowhub/gincontrib/jwt"
	"fmt"
	"encoding/base64"
	"github.com/itsjamie/gin-cors"
	"time"
	"os"
)

var mongoDb db.MongodbAppMetaStore
var fileStore db.FileStore
var appStore store.AppStore
var iotMsgStore *store.IotMsgStore
var appRestController controller.AppRestController
var iotMsgRestController controller.IotMsgRestController

func initAppStore(dbConn string , dbName string , fsLocation string) (error){
	err := mongoDb.Connect(dbConn,dbName)
	if err==nil {
		fileStore = db.FileStore{fsLocation}
		appStore = store.AppStore{&mongoDb, &fileStore}
		iotMsgStore = store.NewIotMsgStore(mongoDb.GetDbConnection())

		appRestController = controller.AppRestController{&appStore}
		iotMsgRestController = controller.IotMsgRestController{iotMsgStore}
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
	r.Static("/bfhub/static","./static")
	r.LoadHTMLGlob("templates/**/*")
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

	root := r.Group("/bfhub/api/apps")
	{
		// public and logged in users
		root.GET("",appRestController.GetApps)
		root.GET("/id/:app_id",appRestController.GetAppById)
		root.GET("/by_full_name",appRestController.GetAppByFullName)
		root.GET("/file/:app_id",appRestController.GetFile)
		root.POST("",appRestController.PostApp)
		root.DELETE("/:app_id",appRestController.DeleteApps)
		//curl -i -v -F "ID=test_file" -F "file=@package.json" http://localhost:8080/bfhub/api/file
		root.POST("/file",appRestController.PostFile)
	}
	root2 := r.Group("/bfhub/api/iotmsg")
	{
		root2.GET("",iotMsgRestController.GetIotMsgs)
		root2.GET("/template",iotMsgRestController.GetIotMsgTemplate)
		root2.POST("",iotMsgRestController.PostIotMsg)
	}
	rootIotMsgUi := r.Group("/bfhub/ui/iotmsg")
	{
		rootIotMsgUi.GET("/list",iotMsgRestController.GetIotMsgsUi)
		rootIotMsgUi.GET("/msg/:msg_id",iotMsgRestController.GetIotMsgUi)

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