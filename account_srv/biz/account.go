package biz

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"gorm.io/gorm"
	"mic-study/account_srv/model"
	"mic-study/account_srv/proto/pb"
	"mic-study/custom_error"
	"mic-study/internal"
)

type AccountServer struct {
	pb.UnimplementedAccountServiceServer
}

func Paginate(pageNo, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageNo == 0 {
			pageNo = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (pageNo - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (a *AccountServer) GetAccountList(ctx context.Context, req *pb.PagingRequest) (*pb.AccountListRes, error) {
	var accountList []model.Account
	result := internal.DB.Scopes(Paginate(int(req.PageNo), int(req.PageSize))).Find(&accountList)
	if result.Error != nil {
		return nil, result.Error
	}
	accountListRes := &pb.AccountListRes{}
	accountListRes.Total = int32(result.RowsAffected)
	for _, account := range accountList {
		accountRes := Model2Pb(account)
		accountListRes.AccountList = append(accountListRes.AccountList, accountRes)
	}
	fmt.Println("GetAccountList Invoked...")
	return accountListRes, nil
}

func Model2Pb(account model.Account) *pb.AccountRes {
	accountRes := &pb.AccountRes{
		Id:       int32(account.ID),
		Mobile:   account.Mobile,
		Password: account.Password,
		Nickname: account.NickName,
		Gender:   account.Gender,
		Role:     uint32(account.Role),
	}
	return accountRes
}

func (a *AccountServer) GetAccountByMobile(ctx context.Context, req *pb.MobileRequest) (*pb.AccountRes, error) {
	var account model.Account
	result := internal.DB.Where(&model.Account{Mobile: req.Mobile}).First(&account)
	if result.RowsAffected == 0 {
		return nil, errors.New(custom_error.AccountNotFound)
	}
	return Model2Pb(account), nil
}

func (a *AccountServer) GetAccountById(ctx context.Context, req *pb.IdRequest) (*pb.AccountRes, error) {
	var account model.Account
	result := internal.DB.First(&account, req.Id)
	if result.RowsAffected == 0 {
		return nil, errors.New(custom_error.AccountNotFound)
	}
	return Model2Pb(account), nil
}

func (a *AccountServer) AddAccount(ctx context.Context, req *pb.AddAccountRequest) (*pb.AccountRes, error) {
	var account model.Account
	result := internal.DB.Where(&model.Account{Mobile: req.Mobile}).First(&account)
	if result.RowsAffected == 1 {
		return nil, errors.New(custom_error.AccountExists)
	}
	account.Mobile = req.Mobile
	account.NickName = req.NickName
	account.Role = 1
	options := password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: md5.New,
	}
	salt, encodePwd := password.Encode(req.Password, &options)
	account.Salt = salt
	account.Password = encodePwd
	r := internal.DB.Create(&account)
	if r.Error != nil {
		return nil, errors.New(custom_error.InternalError)
	}
	return Model2Pb(account), nil
}

func (a *AccountServer) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.UpdateAccountRes, error) {
	var account model.Account
	res := internal.DB.First(&account, req.Id)
	if res.RowsAffected == 0 {
		return nil, errors.New(custom_error.AccountNotFound)
	}
	account.Mobile = req.Mobile
	account.NickName = req.NickName
	account.Gender = req.Gender
	res = internal.DB.Save(&account)
	if res.Error != nil {
		return nil, errors.New(custom_error.InternalError)
	}
	return &pb.UpdateAccountRes{Result: true}, nil
}

func (a *AccountServer) CheckPassword(ctx context.Context, req *pb.CheckPasswordRequest) (*pb.CheckPasswordRes, error) {
	var account model.Account
	res := internal.DB.First(&account, req.AccountId)
	if res.Error != nil {
		return nil, errors.New(custom_error.InternalError)
	}
	if account.Salt == "" {
		return nil, errors.New(custom_error.SaltError)
	}
	options := password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: md5.New,
	}
	r := password.Verify(req.Password, account.Salt, account.Password, &options)
	return &pb.CheckPasswordRes{Result: r}, nil
}
