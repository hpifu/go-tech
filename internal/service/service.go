package service

import (
	"github.com/hpifu/go-account/pkg/account"
	godtoken "github.com/hpifu/go-godtoken/api"
	"github.com/hpifu/go-tech/internal/mysql"
	"github.com/sirupsen/logrus"
)

var InfoLog *logrus.Logger = logrus.New()
var WarnLog *logrus.Logger = logrus.New()
var AccessLog *logrus.Logger = logrus.New()

type Service struct {
	db          *mysql.Mysql
	accountCli  *account.Client
	godtokenCli *godtoken.ServiceClient
}

func NewService(
	db *mysql.Mysql,
	accountCli *account.Client,
	godtokenCli *godtoken.ServiceClient,
) *Service {
	return &Service{
		db:          db,
		accountCli:  accountCli,
		godtokenCli: godtokenCli,
	}
}
