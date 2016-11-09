package db

import (
	"os"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"strings"

	"github.com/unicok/auth/proto/auth"
	db "github.com/unicok/misc/db/mongodb"
	"github.com/unicok/misc/log"
)

var (
	Url    = "mongodb://username:password@127.0.0.1:27017/account"
	dbName string
	accDB  *db.DialContext
)

func Init() {
	var err error
	accDB, err = db.Dial(Url, db.DefaultConcurrent)
	if err != nil {
		log.Panicf("mongodb: cannot connect to %v, err: %v", Url, err)
		os.Exit(-1)
	}

	sl := strings.Split(Url, "/")
	dbName = sl[len(sl)-1]
	if len(dbName) == 0 {
		log.Panic("mongodb: url need db name")
		os.Exit(-1)
	}

	log.Info("connected mongodb:", Url)
}

func Create(user *auth.User) error {

	var newAcc = Account{
		Id:         bson.NewObjectId(),
		UserId:     user.UserId,
		UserName:   user.UserName,
		Password:   user.Password,
		DeviceName: user.DeviceName,
		DeviceId:   user.DeviceId,
		OpenUUID:   user.OpenUUID,
		Lang:       user.Lang,
		LoginIP:    user.LoginIP,
		Created:    time.Now(),
		Updated:    time.Now(),
	}

	err := DoAction(func(c *mgo.Collection) error {
		return c.Insert(&newAcc)
	})

	return err
}

func FindByName(username string) (*Account, error) {
	acc := Account{}
	err := DoAction(func(c *mgo.Collection) error {
		return c.Find(bson.M{"username": username}).One(&acc)
	})
	return &acc, err
}

func DoAction(f db.DBActionFunc) error {
	return accDB.DBAction(dbName, "account", f)
}
