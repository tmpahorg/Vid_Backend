package service

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xproperty"
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/database/helper"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/provide/sn"
)

type VideoService struct {
	db          *gorm.DB
	userService *UserService

	_orderByFunc func(string) string
}

func NewVideoService() *VideoService {
	return &VideoService{
		db:           xdi.GetByNameForce(sn.SGorm).(*gorm.DB),
		userService:  xdi.GetByNameForce(sn.SUserService).(*UserService),
		_orderByFunc: xproperty.GetMapperDefault(&dto.VideoDto{}, &po.Video{}).ApplyOrderBy,
	}
}

func (v *VideoService) QueryAll(pageOrder *param.PageOrderParam) (videos []*po.Video, total int32) {
	total = 0
	v.db.Model(&po.Video{}).Count(&total)

	videos = make([]*po.Video, 0)
	helper.GormPager(v.db, pageOrder.Limit, pageOrder.Page).
		Model(&po.Video{}).
		Order(v._orderByFunc(pageOrder.Order)).
		Find(&videos)

	for _, video := range videos {
		user := &po.User{}
		rdb := v.db.Model(&po.User{}).Where(&po.User{Uid: video.AuthorUid}).First(user)
		if !rdb.RecordNotFound() {
			video.Author = user
		}
	}

	return videos, total
}

func (v *VideoService) QueryByUid(uid int32, pageOrder *param.PageOrderParam) (videos []*po.Video, total int32, status database.DbStatus) {
	author := v.userService.QueryByUid(uid)
	if author == nil {
		return nil, 0, database.DbNotFound
	}

	total = 0
	v.db.Model(&po.Video{}).Where(&po.Video{AuthorUid: uid}).Count(&total)

	videos = make([]*po.Video, 0)
	helper.GormPager(v.db, pageOrder.Limit, pageOrder.Page).
		Model(&po.Video{}).
		Order(v._orderByFunc(pageOrder.Order)).
		Where(&po.Video{AuthorUid: uid}).
		Find(&videos)

	for idx := range videos {
		videos[idx].Author = author
	}

	return videos, total, database.DbSuccess
}

func (v *VideoService) QueryCountByUid(uid int32) (int32, database.DbStatus) {
	if !v.userService.Exist(uid) {
		return 0, database.DbNotFound
	}

	total := 0
	v.db.Model(&po.Video{}).Where(&po.Video{AuthorUid: uid}).Count(&total)

	return int32(total), database.DbSuccess
}

func (v *VideoService) QueryByVid(vid int32) *po.Video {
	video := &po.Video{}
	rdb := v.db.Model(&po.Video{}).Where(&po.Video{Vid: vid}).First(video)
	if rdb.RecordNotFound() {
		return nil
	}

	user := &po.User{}
	rdb = v.db.Model(&po.User{}).Where(&po.User{Uid: video.AuthorUid}).First(&user)
	if !rdb.RecordNotFound() {
		video.Author = user
	}

	return video
}

func (v *VideoService) Exist(vid int32) bool {
	return helper.GormExist(v.db, &po.Video{}, &po.Video{Vid: vid})
}

func (v *VideoService) Insert(video *po.Video) database.DbStatus {
	return helper.GormInsert(v.db, &po.Video{}, video)
}

func (v *VideoService) Update(video *po.Video) database.DbStatus {
	return helper.GormUpdate(v.db, &po.Video{}, video)
}

func (v *VideoService) Delete(vid int32) database.DbStatus {
	return helper.GormDelete(v.db, &po.Video{}, &po.Video{Vid: vid})
}

func (v *VideoService) DeleteBy2Id(vid int32, uid int32) database.DbStatus {
	return helper.GormDelete(v.db, &po.Video{}, &po.Video{Vid: vid, AuthorUid: uid})
}
