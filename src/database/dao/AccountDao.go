package dao

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/model/po"
)

type AccountDao struct {
	Db *gorm.DB `di:"~"`
}

func NewPassDao(dic *xdi.DiContainer) *AccountDao {
	repo := &AccountDao{}
	if !dic.Inject(repo) {
		panic("Inject failed")
	}
	return repo
}

func (a *AccountDao) QueryByUsername(username string) *po.Account {
	user := QueryHelper(a.Db, &po.User{}, &po.User{Username: username}).(*po.User)
	if user == nil {
		return nil
	}
	account := QueryHelper(a.Db, &po.Account{}, &po.Account{Uid: user.Uid}).(*po.Account)
	if account == nil {
		return nil
	}
	account.User = user
	return account
}

func (a *AccountDao) Insert(pass *po.Account) database.DbStatus {
	return InsertHelper(a.Db, &po.Account{}, pass) // cascade create
}

func (a *AccountDao) Update(pass *po.Account) database.DbStatus {
	return UpdateHelper(a.Db, &po.Account{}, pass)
}
