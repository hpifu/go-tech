package service

import (
	"fmt"

	"github.com/hpifu/go-kit/hhttp"
	"github.com/hpifu/go-tech/internal/mysql"
	"github.com/sirupsen/logrus"
)

var InfoLog *logrus.Logger = logrus.New()
var WarnLog *logrus.Logger = logrus.New()
var AccessLog *logrus.Logger = logrus.New()

type Service struct {
	db         *mysql.Mysql
	client     *hhttp.HttpClient
	apiAccount string
}

func NewService(
	db *mysql.Mysql,
	client *hhttp.HttpClient,
	apiAccount string,
) *Service {
	return &Service{
		db:         db,
		client:     client,
		apiAccount: apiAccount,
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

func (s *Service) getAccount(token string) (*Account, error) {
	result := s.client.GET("http://"+s.apiAccount+"/account/"+token, nil, nil)
	if result.Err != nil {
		return nil, result.Err
	}

	if result.Status == 204 {
		return nil, nil
	}

	if result.Status == 200 {
		res := &Account{}
		if err := result.Interface(res); err != nil {
			return nil, err
		}

		return res, nil
	}

	return nil, fmt.Errorf("GET account failed. res[%v]", string(result.Res))
}
