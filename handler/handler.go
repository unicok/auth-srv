package handler

import (
	"errors"
	"os"
	"regexp"
	"unsafe"

	mgo "gopkg.in/mgo.v2"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"github.com/unicok/auth/db"
	proto "github.com/unicok/auth/proto/auth"
	"github.com/unicok/misc/log"
	"github.com/unicok/snowflake/proto/snowflake"
	"golang.org/x/net/context"
)

var (
	ErrMethodNotSupported      = errors.New("method not supported")
	ErrPasswordInvalid         = errors.New("password invalid")
	ErrUsernameOrPasswordEmpty = errors.New("username or password empty")
	ErrUsernameAlreadyExists   = errors.New("username already exists")

	AuthFailResult = &proto.Result{OK: false, UserId: 0, Body: nil}

	uuidRegexp = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
)

var (
	sfCli snowflake.SnowflakeServiceClient
)

type Auth struct{}

func (s *Auth) Register(ctx context.Context, req *proto.RegisterRequest, rsp *proto.Result) error {
	if len(req.User.UserName) == 0 || len(req.User.Password) == 0 {
		rsp = AuthFailResult
		return ErrUsernameOrPasswordEmpty
	}

	_, err := db.FindByName(req.User.UserName)
	if err == nil {
		rsp = AuthFailResult
		return ErrUsernameAlreadyExists
	}

	if err != nil && err != mgo.ErrNotFound {
		rsp.OK = false
		return err
	}

	// 创建ID
	req.User.UserId, err = GetNextID("userid")
	if err != nil {
		rsp = AuthFailResult
		return err
	}

	err = db.Create(req.User)
	if err != nil {
		rsp = AuthFailResult
		return err
	}

	rsp = &proto.Result{OK: true, UserId: 0, Body: nil}
	return nil
}

func (s *Auth) Auth(ctx context.Context, req *proto.Certificate, rsp *proto.Result) error {
	return nil
}

func GetNextID(key string) (uint64, error) {
	req := client.NewRequest(
		"com.unicok.srv.snowflake", "SnowflakeService.Next",
		&snowflake.Snowflake_Key{
			Name: "userid",
		})

	// create context with metadata
	ctx := metadata.NewContext(context.Background(), nil)

	rsp := &snowflake.Snowflake_Value{}
	if err := client.Call(ctx, req, rsp); err != nil {
		return 0, err
	}

	log.Info("new userid:", rsp.Value)
	return uint64(rsp.Value), nil
}

func checkErrPanic(err error) {
	if err != nil {
		log.Panic(err)
		os.Exit(-1)
	}
}

func Str2Bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
