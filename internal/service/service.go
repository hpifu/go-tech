package service

import (
	"github.com/hpifu/go-account/pkg/account"
	"github.com/hpifu/go-tech/internal/mysql"
	"github.com/sirupsen/logrus"
)

var InfoLog *logrus.Logger = logrus.New()
var WarnLog *logrus.Logger = logrus.New()
var AccessLog *logrus.Logger = logrus.New()

type Service struct {
	db     *mysql.Mysql
	client *account.Client
}

func NewService(
	db *mysql.Mysql,
	client *account.Client,
) *Service {
	return &Service{
		db:     db,
		client: client,
	}
}

type Account struct {
	ID        int    `form:"id" json:"id,omitempty"`
	Email     string `form:"email" json:"email,omitempty"`
	Phone     string `form:"phone" json:"phone,omitempty"`
	FirstName string `form:"firstName" json:"firstName,omitempty"`
	LastName  string `form:"lastName" json:"lastName,omitempty"`
	Birthday  string `form:"birthday" json:"birthday,omitempty"`
	Password  string `form:"password" json:"password,omitempty"`
	Gender    int    `form:"gender" json:"gender"`
	Avatar    string `form:"avatar" json:"avatar"`
}
