package service

import (
	"context"
	"fmt"
	"github.com/hpifu/go-account/pkg/account"
	godtoken "github.com/hpifu/go-godtoken/api"
	"github.com/hpifu/go-tech/internal/mysql"
	"github.com/sirupsen/logrus"
	"time"
)

var InfoLog *logrus.Logger = logrus.New()
var WarnLog *logrus.Logger = logrus.New()
var AccessLog *logrus.Logger = logrus.New()

type Service struct {
	db              *mysql.Mysql
	accountCli      *account.Client
	godtokenCli     godtoken.ServiceClient
	godtokenTimeout time.Duration
}

func NewService(
	db *mysql.Mysql,
	accountCli *account.Client,
	godtokenCli godtoken.ServiceClient,
) *Service {
	return &Service{
		db:              db,
		accountCli:      accountCli,
		godtokenCli:     godtokenCli,
		godtokenTimeout: 200 * time.Millisecond,
	}
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
