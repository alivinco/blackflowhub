package controller

import (
	"github.com/alivinco/blackflowhub/store"
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/alivinco/blackflowhub/model"
	"github.com/alivinco/blackflowhub/gincontrib/utils"
	"net/http"
	"time"
	"gopkg.in/mgo.v2/bson"
)

type IotMsgRestController struct {
	IotMsgStore *store.IotMsgStore
}

func(ac *IotMsgRestController) PostIotMsg(c *gin.Context){
	auth := utils.GetAuthRequest(c)
	// Only authenticated requests are allowed
	if auth.IsAuthenticated{
		var iotMsgModel model.IotMsg
		err := c.BindJSON(&iotMsgModel)
		if  err== nil {
			// Check if developer in model is the same as in Claim
			if iotMsgModel.Author == auth.Username {
				if iotMsgModel.ID.Hex() == ""{
					iotMsgModel.Created = time.Now()
				}
				iotMsgModel.Updated = time.Now()
				ID, err := ac.IotMsgStore.UpsertMsg(&iotMsgModel)
				if err == nil {
					fmt.Println(iotMsgModel)
					c.JSON(http.StatusOK, gin.H{"status": "OK", "app_id":ID})
				}else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				}
			}else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not allowed to access the msg"})
			}

		}else {
		   c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't bind IotMsg model.Error"+err.Error()})
		}
	}else{
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not allowed to access the msg"})
	}

}

func(ac *IotMsgRestController) GetIotMsgTemplate(c *gin.Context) {
	result := model.IotMsg{}
	result.Created = time.Now()
	result.Updated = time.Now()
	comments := make([]model.Comment,1)
	comments[0] = model.Comment{}
	result.Comments = comments
	c.JSON(http.StatusOK,result)
}


func(ac *IotMsgRestController) GetIotMsgs(c *gin.Context) {
	auth := utils.GetAuthRequest(c)
	var filter model.IotMsg
	if  auth.IsAuthenticated {
		fmt.Println("Sending back only private message fow the user :",auth.Username)
		filter = model.IotMsg{Author:auth.Username}
	}else {
		fmt.Println("Request is not authenticated , sending only public apps")
		filter = model.IotMsg{Public:true}
	}

	result , err := ac.IotMsgStore.GetIotMsgs(&filter)
	if err == nil {
		c.JSON(http.StatusOK,*result)
	} else {
		c.JSON(http.StatusInternalServerError,gin.H{"error": err.Error()})
	}
}

func(ac *IotMsgRestController) GetIotMsgsUi(c *gin.Context) {
	auth := utils.GetAuthRequest(c)
	var filter model.IotMsg
	if  auth.IsAuthenticated {
		fmt.Println("Sending back only private message fow the user :",auth.Username)
		filter = model.IotMsg{Author:auth.Username}
	}else {
		fmt.Println("Request is not authenticated , sending only public apps")
		filter = model.IotMsg{Public:true}
	}

	result , err := ac.IotMsgStore.GetIotMsgs(&filter)
	if err == nil {
		c.HTML(http.StatusOK,"iot_msg_list.html",*result)
	} else {
		c.JSON(http.StatusInternalServerError,gin.H{"error": err.Error()})
	}
}

func(ac *IotMsgRestController) GetIotMsgUi(c *gin.Context) {
	ID := c.Param("msg_id")
	auth := utils.GetAuthRequest(c)
	var filter model.IotMsg
	if  auth.IsAuthenticated {
		fmt.Println("Sending back only private message fow the user :",auth.Username)
		filter = model.IotMsg{Author:auth.Username}
	}else {
		fmt.Println("Request is not authenticated , sending only public apps")
		filter = model.IotMsg{Public:true}
	}
	filter.ID = bson.ObjectIdHex(ID)
	result , err := ac.IotMsgStore.GetIotMsgs(&filter)
	if err == nil {
		if len(*result)==1 {
			c.HTML(http.StatusOK,"iot_msg.html",(*result)[0])
		}else{
			c.HTML(http.StatusNoContent,"iot_msg.html",gin.H{})
		}

	} else {
		c.JSON(http.StatusInternalServerError,gin.H{"error": err.Error()})
	}
}

