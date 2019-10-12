package service

import (
	"github.com/hpifu/go-tech/internal/mysql"
	"github.com/sirupsen/logrus"
)

var InfoLog *logrus.Logger
var WarnLog *logrus.Logger
var AccessLog *logrus.Logger

func init() {
	InfoLog = logrus.New()
	WarnLog = logrus.New()
	AccessLog = logrus.New()
}

type Service struct {
	secure bool
	domain string
	db     *mysql.Mysql
}

func NewService(
	secure bool,
	domain string,
	db *mysql.Mysql,
) *Service {
	return &Service{
		secure: secure,
		domain: domain,
		db:     db,
	}
}
