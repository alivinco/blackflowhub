package model
import (
	"gopkg.in/mgo.v2/bson"
	"time"
)


type IotMsg struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"id"`
	PayloadType string `json:"payload_type"`
	Type string `json:"type"`
	Class string `json:"class"`
	SubClass string `json:"sub_class"`
	Version string `json:"version"`
	Author string `json:"author"`
	Tags []string `json:"tags"`
	Categories []string `json:"categories"`
	Description string `json:"description"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Comments []Comment `json:"comments"`
	TopicTemplates []string `json:"topic_templates"`
	ReviewStatus string `json:"review_status"`
	Public bool `json:"public"`
	MsgExamples []string `json:"msg_examples"`
}

type Comment struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Author string `json:"author"`
	Created time.Time `json:"created"`
	Message string `json:"message"`
}