package model
import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type App struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"id"`
	AppName string `json:"app_name"`
	Version float32 `json:"version"`
	Tags []string `json:"tags"`
	Categories []string `json:"categories"`
	ShortDescription string `json:"short_description"`
	LongDescription string `json:"long_description"`
	Rating int `json:"rating"`
	Downloads int `json:"downloads"`
	Meta map[string]string `json:"meta"`
	Technology string `json:"technology"`
	File string `json:"file"`
	Developer string `json:"developer"`
	Size int32 `json:"size"`
	Compatibility string `json:"compatibility"`
	// timestamp as UCT time in seconds .
	Updated time.Time `json:"updated"`
	Public bool `json:"public"`
}

type AppFilter struct {
	AppName string `json:"app_name"`
	Version float32 `json:"version"`
	Tags []string `json:"tags"`
	Categories string `json:"categories"`
	Developer string `json:"developer"`
	Public bool `json:"public"`
}