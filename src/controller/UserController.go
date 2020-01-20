package controller

import (
	"github.com/Aoi-hosizora/ahlib/xcondition"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xmapper"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/controller/exception"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/database/dao"
	"github.com/vidorg/vid_backend/src/middleware"
	"github.com/vidorg/vid_backend/src/model/common"
	"github.com/vidorg/vid_backend/src/model/common/enum"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/dto/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/util"
	"net/http"
)

type UserController struct {
	Config     *config.ServerConfig   `di:"~"`
	JwtService *middleware.JwtService `di:"~"`
	UserDao    *dao.UserDao           `di:"~"`
	VideoDao   *dao.VideoDao          `di:"~"`
	SubDao     *dao.SubDao            `di:"~"`
	Mapper     *xmapper.EntityMapper  `di:"~"`
}

func NewUserController(dic *xdi.DiContainer) *UserController {
	ctrl := &UserController{}
	dic.Inject(ctrl)
	if xdi.HasNilDi(ctrl) {
		panic("Has nil di field")
	}

	return ctrl
}

// @Router				/v1/user?page [GET] [Auth]
// @Summary				查询所有用户
// @Description			管理员查询所有用户，返回分页数据，管理员权限，此处可见用户手机号码
// @Tag					User
// @Tag					Administration
// @Param				page query integer false "分页"
// @Accept				multipart/form-data
// @ErrorCode			400 request param error
// @ErrorCode			401 need admin authority
/* @Response 200		{
							"code": 200,
							"message": "success",
							"data": {
								"count": 1,
								"page": 1,
								"data": [ ${user} ]
							}
 						} */
func (u *UserController) QueryAllUsers(c *gin.Context) {
	page, ok := param.BindQueryPage(c)
	if !ok {
		common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
		return
	}

	users, count := u.UserDao.QueryAll(page)

	// show all user's info
	retDto := xcondition.First(u.Mapper.Map([]*dto.UserDto{}, users, dto.UserDtoAdminMapOption())).([]*dto.UserDto)
	common.Result{}.Ok().SetPage(count, page, retDto).JSON(c)
}

// @Router				/v1/user/{uid} [GET]
// @Summary				查询用户
// @Description			查询用户信息，此处可见用户手机号码
// @Tag					User
// @Param				uid path integer true "用户id"
// @Accept				multipart/form-data
// @ErrorCode			400 request param error
// @ErrorCode			404 user not found
/* @Response 200		{
							"code": 200,
							"message": "success",
							"data": {
								"user": ${user},
								"extra": {
									"subscribing_cnt": 3,
									"subscriber_cnt": 2,
									"video_cnt": 3
								}
							}
 						} */
func (u *UserController) QueryUser(c *gin.Context) {
	uid, ok := param.BindRouteId(c, "uid")
	if !ok {
		common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
		return
	}

	user := u.UserDao.QueryByUid(uid)
	if user == nil {
		common.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()).JSON(c)
		return
	}
	subscribingCnt, subscriberCnt, _ := u.SubDao.QuerySubCnt(user.Uid)
	videoCnt, _ := u.VideoDao.QueryCount(user.Uid)
	extraInfo := &dto.UserExtraInfo{
		SubscribingCount: subscribingCnt,
		SubscriberCount:  subscriberCnt,
		VideoCount:       videoCnt,
	}

	// need to squeeze out whether you can see the admin info
	authUser := u.JwtService.GetAuthUser(c)
	retDto := xcondition.First(u.Mapper.Map(&dto.UserDto{}, user, dto.UserDtoExtraMapOption(authUser))).(*dto.UserDto)
	common.Result{}.Ok().PutData("user", retDto).PutData("extra", extraInfo).JSON(c)
}

// @Router				/v1/user/ [PUT] [Auth]
// @Summary				更新登录用户
// @Description			更新用户信息
// @Tag					User
// @Param				username formData string true "用户名，长度在 [8, 30] 之间"
// @Param				sex formData string true "用户性别，允许值为 {male, female, unknown}"
// @Param				profile formData string true "用户简介，长度在 [0, 255] 之间"
// @Param				birth_time formData string true "用户生日，固定格式为 2000-01-01"
// @Param				phone_number formData string true "用户手机号码，长度为 11，仅限中国大陆手机号码"
// @Param				avatar_file formData file true "用户头像链接"
// @Accept				multipart/form-data
// @ErrorCode			400 request param error
// @ErrorCode			400 request format error
// @ErrorCode			400 username has been used
// @ErrorCode			404 user not found
// @ErrorCode			500 user update failed
/* @Response 200		{
							"code": 200,
							"message": "success",
							"data": ${user}
 						} */
// @Router				/v1/user/admin/{uid} [PUT] [Auth]
// @Summary				管理员更新用户
// @Description			更新用户信息，管理员权限
// @Tag					User
// @Tag					Administration
// @Param				username formData string true "用户名，长度在 [8, 30] 之间"
// @Param				sex formData string true "用户性别，允许值为 {male, female, unknown}"
// @Param				profile formData string true "用户简介，长度在 [0, 255] 之间"
// @Param				birth_time formData string true "用户生日，固定格式为 2000-01-01"
// @Param				phone_number formData string true "用户手机号码，长度为 11，仅限中国大陆手机号码"
// @Param				avatar_file formData file true "用户头像链接"
// @Accept				multipart/form-data
// @ErrorCode			400 request param error
// @ErrorCode			400 request format error
// @ErrorCode			400 username has been used
// @ErrorCode			401 need admin authority
// @ErrorCode			404 user not found
// @ErrorCode			500 user update failed
/* @Response 200		{
							"code": 200,
							"message": "success",
							"data": ${user}
 						} */
func (u *UserController) UpdateUser(isExact bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := &po.User{}
		if !isExact {
			user = u.JwtService.GetAuthUser(c)
		} else {
			uid, ok := param.BindRouteId(c, "uid")
			if !ok {
				common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
				return
			}
			user = u.UserDao.QueryByUid(uid)
			if user == nil {
				common.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()).JSON(c)
				return
			}
		}
		// Update
		userParam := &param.UserParam{}
		if err := c.ShouldBind(userParam); err != nil {
			common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.WrapValidationError(err).Error()).JSON(c)
			return
		}
		user.Username = userParam.Username
		user.Sex = enum.StringToSexType(userParam.Sex)
		user.Profile = *userParam.Profile
		user.BirthTime, _ = common.JsonDate{}.Parse(userParam.BirthTime)
		user.PhoneNumber = userParam.PhoneNumber
		url, ok := util.CommonUtil.GetFilenameFromUrl(userParam.AvatarUrl, u.Config.FileConfig.ImageUrlPrefix)
		if !ok {
			common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
			return
		}
		user.AvatarUrl = url

		status := u.UserDao.Update(user)
		if status == database.DbNotFound {
			common.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()).JSON(c)
			return
		} else if status == database.DbExisted {
			common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.UsernameUsedError.Error()).JSON(c)
			return
		} else if status == database.DbFailed {
			common.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.UserUpdateError.Error()).JSON(c)
			return
		}

		retDto := xcondition.First(u.Mapper.Map(&dto.UserDto{}, user)).(*dto.UserDto)
		common.Result{}.Ok().SetData(retDto).JSON(c)
	}
}

// @Router				/v1/user/ [DELETE] [Auth]
// @Summary				删除登录用户
// @Description			删除用户所有信息
// @Tag					User
// @Accept				multipart/form-data
// @ErrorCode			404 user not found
// @ErrorCode			500 user delete failed
/* @Response 200		{
							"code": 200,
							"message": "success"
 						} */
// @Router				/v1/user/admin/{uid} [DELETE] [Auth]
// @Summary				管理员删除用户
// @Description			删除用户所有信息，管理员权限
// @Tag					User
// @Tag					Administration
// @Accept				multipart/form-data
// @ErrorCode			401 need admin authority
// @ErrorCode			404 user not found
// @ErrorCode			500 user delete failed
/* @Response 200		{
							"code": 200,
							"message": "success"
 						} */
func (u *UserController) DeleteUser(isExact bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		var uid int32
		if !isExact {
			uid = u.JwtService.GetAuthUser(c).Uid
		} else {
			_uid, ok := param.BindRouteId(c, "uid")
			if !ok {
				common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
				return
			}
			uid = _uid
		}
		// Delete
		status := u.UserDao.Delete(uid)
		if status == database.DbNotFound {
			common.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()).JSON(c)
			return
		} else if status == database.DbFailed {
			common.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.UserDeleteError.Error()).JSON(c)
			return
		}
		common.Result{}.Ok().JSON(c)
	}
}
