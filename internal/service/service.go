package service

import (
	"github.com/hpifu/go-tech/internal/mysql"
	"github.com/sirupsen/logrus"
)

var InfoLog *logrus.Logger = logrus.New()
var WarnLog *logrus.Logger = logrus.New()
var AccessLog *logrus.Logger = logrus.New()

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
