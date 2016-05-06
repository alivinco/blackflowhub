package controller
import (
	"github.com/gin-gonic/gin"
	"github.com/alivinco/blackflowhub/model"
	"net/http"
	"fmt"
	"github.com/alivinco/blackflowhub/store"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"os"
	"io"
	"github.com/alivinco/blackflowhub/gincontrib/utils"
	"errors"
)
type AppRestController struct {
	AppStore *store.AppStore
}

func(ac *AppRestController) PostApp(c *gin.Context){
	auth := utils.GetAuthRequest(c)
	// Only authenticated requests are allowed
	if auth.IsAuthenticated{
		var appModel model.App
		if c.BindJSON(&appModel) == nil {
			// Check if developer in model is the same as in Claim
			if appModel.Developer == auth.Username {
				ID, err := ac.AppStore.UpsertApp(&appModel)
				if err == nil {
					fmt.Println(appModel)
					c.JSON(http.StatusOK, gin.H{"status": "OK", "app_id":ID})
				}else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				}
			}else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not allowed to access the app"})
			}

		}else {
		   c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't bind model"})
		}
	}else{
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not allowed to access the app"})
	}

}

func(ac *AppRestController) DeleteApps(c *gin.Context) {
	ID := c.Param("app_id")
	fmt.Println("Deleting app with id :",ID)
	// Only Developer can delete the app
	if ac.AuthorizeAppUpdate(c,ID) {
		err := ac.AppStore.DeleteApp(ID)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"status": "OK"})
		}else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}else{
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not allowed to access the app"})
	}
}

func(ac *AppRestController) GetApps(c *gin.Context) {
	auth := utils.GetAuthRequest(c)
	var filter model.AppFilter
	if  auth.IsAuthenticated {
		fmt.Println("Sending back only private appf fow user :",auth.Username)
		filter = model.AppFilter{Developer:auth.Username}
	}else {
		fmt.Println("Request is not authenticated , sending only public apps")
		filter = model.AppFilter{Public:true}
	}

	result , err := ac.AppStore.GetApps(&filter)
	if err == nil {
		c.JSON(http.StatusOK,*result)
	} else {
		c.JSON(http.StatusInternalServerError,gin.H{"error": err.Error()})
	}
}

func(ac *AppRestController) GetAppById(c *gin.Context) {
	auth := utils.GetAuthRequest(c)
	result , err := ac.AppStore.GetAppById(c.Param("app_id"))
	if err == nil {
		if auth.Username == result.Developer || result.Public {
			c.JSON(http.StatusOK, *result)
		}else{
			c.AbortWithError(401, errors.New("You are not allowed to access the app"))
		}

	} else {
		c.JSON(http.StatusNotFound,gin.H{"error": err.Error()})
	}
}

func(ac *AppRestController) GetAppByFullName(c *gin.Context) {
	auth := utils.GetAuthRequest(c)
	version , err := strconv.ParseFloat(c.Query("version"),32)
	if err == nil {
		result, err := ac.AppStore.GetAppByFullName(c.Query("developer"), c.Query("app_name"), float32(version))
		if err == nil {
			    if auth.Username == result.Developer || result.Public {
					c.JSON(http.StatusOK, *result)
				}else{
					c.AbortWithError(401, errors.New("You are not allowed to access the app"))
				}
			}else{
				c.JSON(http.StatusNotFound,gin.H{"status":"not found"})
		     }

	}else
	{
		c.JSON(http.StatusInternalServerError,gin.H{"error": err.Error()})
	}
}

//func(ac *AppRestController) GetAppByName(c *gin.Context) {
//	auth := utils.GetAuthRequest(c)
//	result , err := ac.AppStore.GetApps(nil)
//	if err == nil {
//		if auth.Username == result.Developer || result.Public {
//			c.JSON(http.StatusOK,*result)
//		c.AbortWithError(401, errors.New("You are not allowed to access the app"))
//	} else {
//		c.JSON(http.StatusInternalServerError,gin.H{"error": err.Error()})
//	}
//}


func(ac *AppRestController) PostFile(c *gin.Context) {
   auth := utils.GetAuthRequest(c)
   if auth.IsAuthenticated {
	   file, header, err := c.Request.FormFile("file")
	   if err == nil {
		   ID := c.PostForm("id")
		   fmt.Println("ID:", ID)
		   filename := header.Filename
		   fmt.Println("File:", filename)
		   // only developer is allowed to modify app

		   if ac.AuthorizeAppUpdate(c,ID) {
			   fileContent, _ := ioutil.ReadAll(file)
			   err = ac.AppStore.SaveFile(ID, filename, fileContent)
			   if err == nil {
				   c.JSON(http.StatusOK, gin.H{"status": "OK"})
			   }else {
				   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			   }
		   }

	   } else {
		   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	   }
   }else {
	   c.AbortWithError(401, errors.New("You are authorized to access the app"))
   }

}
func(ac *AppRestController) GetFile(c *gin.Context) {
	auth := utils.GetAuthRequest(c)
	ID := c.Param("app_id")
	fileFullPath := filepath.Join(ac.AppStore.AppFileStore.GetStorageLocation(),ID)
	am , err := ac.AppStore.AppMetaStore.GetById(ID)
	if err == nil {
		if auth.Username == am.Developer || am.Public {
			am.Downloads = am.Downloads + 1
			ac.AppStore.AppMetaStore.UpdateById(ID, am)
			f, _ := os.Open(fileFullPath)
			defer f.Close()
			c.Header("Content-Type", "application/octet-stream")
			c.Header("Content-Disposition", "attachment; filename=" + am.File)
			io.Copy(c.Writer, f)
		}else{
			c.AbortWithError(401, errors.New("You are authorized to access the app"))
		}
	}else{
		c.JSON(http.StatusNotFound,gin.H{"status":"not found"})
	}
}

func(ac *AppRestController) AuthorizeAppUpdate(c *gin.Context,ID string)(bool){
	// Get Auth data
	auth := utils.GetAuthRequest(c)
	// Get app metadata
	am , err := ac.AppStore.AppMetaStore.GetById(ID)
	if err != nil{
		if am.Developer == auth.Username {
			return true
		}else
		{
			c.AbortWithError(401, errors.New("You are authorized to access the app"))
		}
	}else{
		c.JSON(http.StatusNotFound,gin.H{"status":"app not found"})
	}
	return false
}