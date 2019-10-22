package service

import (
	"context"
	"fmt"
	"github.com/hpifu/go-account/pkg/account"
	godtoken "github.com/hpifu/go-godtoken/api"
	"github.com/hpifu/go-tech/internal/es"
	"github.com/hpifu/go-tech/internal/mysql"
	"github.com/sirupsen/logrus"
	"time"
)

type Service struct {
	db              *mysql.Mysql
	es              *es.ES
	accountCli      *account.Client
	godtokenCli     godtoken.ServiceClient
	godtokenTimeout time.Duration
	infoLog         *logrus.Logger
	warnLog         *logrus.Logger
	accessLog       *logrus.Logger
}

func NewService(
	db *mysql.Mysql,
	es *es.ES,
	accountCli *account.Client,
	godtokenCli godtoken.ServiceClient,
) *Service {
	return &Service{
		db:              db,
		es:              es,
		accountCli:      accountCli,
		godtokenCli:     godtokenCli,
		godtokenTimeout: 200 * time.Millisecond,
		infoLog:         logrus.New(),
		warnLog:         logrus.New(),
		accessLog:       logrus.New(),
	}
}

func (s *Service) SetLogger(infoLog, warnLog, accessLog *logrus.Logger) {
	s.infoLog = infoLog
	s.warnLog = warnLog
	s.accessLog = accessLog
}

func (s *Service) GetAccounts(rid string, ids []int) (map[int]*account.Account, error) {
	var idUnique []int
	idMap := map[int]struct{}{}
	for _, id := range ids {
		idMap[id] = struct{}{}
	}
	for k := range idMap {
		idUnique = append(idUnique, k)
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.godtokenTimeout)
	defer cancel()
	res, err := s.godtokenCli.GetToken(ctx, &godtoken.GetTokenReq{Rid: rid})
	if err != nil {
		return nil, fmt.Errorf("godtoken verify failed. err: [%v]", err)
	}

	accounts, err := s.accountCli.GETAccounts(rid, res.Token, idUnique)
	if err != nil {
		return nil, fmt.Errorf("get accounts failed. err: [%v]", err)
	}
	accountMap := map[int]*account.Account{}
	if accounts != nil {
		for _, a := range accounts {
			accountMap[a.ID] = a
		}
	}

	return accountMap, nil
}

func runecut(content string, length int) string {
	runes := []rune(content)
	var rs []rune
	for _, r := range runes {
		if _, ok := map[rune]struct{}{
			'#': {}, '`': {}, ' ': {}, '\n': {}, '\r': {},
		}[r]; !ok {
			rs = append(rs, r)
		}
	}
	if len(rs) >= length {
		return string(rs[0:length])
	}

	return string(rs)
}
