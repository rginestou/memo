package main

import (
	"html/template"
	"log"
	"os"
	"strings"
	"time"

	"github.com/russross/blackfriday"
	"gopkg.in/mgo.v2/bson"

	mgo "gopkg.in/mgo.v2"
)

// Memo model
type Memo struct {
	ID        bson.ObjectId `bson:"_id" json:"_id"`
	Name      string        `bson:"name" json:"name"`
	Category  string        `bson:"category" json:"category"`
	Content   template.HTML `bson:"content" json:"content"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}

// ToID ...
func (c Memo) ToID() string {
	return c.ID.Hex()
}

// ToMarkupContent returns the HTML from the markdown content
func (c Memo) ToMarkupContent() template.HTML {
	body := string(blackfriday.Run([]byte(c.Content)))
	body = strings.Replace(body, "\n", "<br>", -1)
	return template.HTML(body)
}

// Avatar profile picture
func (c Memo) Avatar() string {
	avatar := "static/avatar.png"
	if _, err := os.Stat("data/" + c.ID.Hex() + ".png"); !os.IsNotExist(err) {
		avatar = "data/" + c.ID.Hex() + ".png"
	}

	return avatar
}

var db *mgo.Database

// Connect to the database
func connect(server, database string) {
	session, err := mgo.Dial(server)
	if err != nil {
		log.Fatal(err)
	}

	db = session.DB(database)
}

func getAllMemos(filter string) ([]Memo, error) {
	var resp []Memo
	err := db.C("memos").Find(bson.M{"$or": []bson.M{
		bson.M{"name": bson.M{"$regex": bson.RegEx{Pattern: "(?i).*" + filter + ".*", Options: "i"}}},
		bson.M{"category": bson.M{"$regex": bson.RegEx{Pattern: "(?i).*" + filter + ".*", Options: "i"}}},
	}}).Sort("name").All(&resp)

	return resp, err
}

func getMemoByID(ID bson.ObjectId) (Memo, error) {
	var resp Memo
	err := db.C("memos").FindId(ID).One(&resp)

	return resp, err
}

func updateMemo(memo Memo) error {
	_, err := db.C("memos").Upsert(bson.M{"_id": memo.ID}, memo)

	return err
}

func deleteMemoByID(ID bson.ObjectId) error {
	return db.C("memos").RemoveId(ID)
}
