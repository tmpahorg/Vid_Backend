package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shomali11/util/xconditions"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/controller/exception"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/database/dao"
	"github.com/vidorg/vid_backend/src/middleware"
	"github.com/vidorg/vid_backend/src/model"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/dto/common"
	"github.com/vidorg/vid_backend/src/model/dto/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/util"
	"net/http"
	"strconv"
	"time"
)

type videoController struct {
	config   *config.ServerConfig
	videoDao *dao.VideoDao
}

func VideoController(config *config.ServerConfig) *videoController {
	return &videoController{
		config:   config,
		videoDao: dao.VideoRepository(config.DatabaseConfig),
	}
}

// @Router 				/video?page [GET] [Auth]
// @Summary 			查询所有视频
// @Description 		管理员查询所有视频，返回分页数据，Admin
// @Tag					Video
// @Tag					Administration
// @Param 				page query integer false "分页"
// @Accept 				multipart/form-data
// @ErrorCode			400 request param error
// @ErrorCode 			401 need admin authority
/* @Success 200 		{
							"code": 200,
							"message": "success",
							"data": {
								"count": 1,
								"page": 1,
								"data": [
									{
										"vid": 10,
										"title": "Video Title",
										"description": "Video Demo Description",
										"video_url": "",
										"cover_url": "http://localhost:3344/raw/image/default/cover.jpg",
										"upload_time": "2019-12-26 14:14:04",
										"update_time": "2019-12-30 21:04:51",
										"author": {
											"uid": 10,
											"username": "aoihosizora",
											"sex": "male",
											"profile": "Demo Profile",
                    						"avatar_url": "http://localhost:3344/raw/image/10/avatar_20191230205334869693.jpg",
											"birth_time": "2019-12-26",
											"authority": "admin"
										}
									}
								]
							}
 						} */
func (v *videoController) QueryAllVideos(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()))
		return
	}
	page = xconditions.IfThenElse(page < 1, 1, page).(int)

	videos, count := v.videoDao.QueryAll(page)
	c.JSON(http.StatusOK, common.Result{}.Ok().SetPage(count, page, dto.VideoDto{}.FromPos(videos)))
}

// @Router 				/user/{uid}/video?page [GET]
// @Summary 			查询用户视频
// @Description 		查询作者为用户的所有视频，返回分页数据
// @Tag					Video
// @Param 				uid path integer true "用户id"
// @Param 				page query integer false "分页"
// @Accept 				multipart/form-data
// @ErrorCode			400 request param error
// @ErrorCode			404 user not found
/* @Success 200 		{
							"code": 200,
							"message": "success",
							"data": {
								"count": 1,
								"page": 1,
								"data": [
									{
										"vid": 10,
										"title": "Video Title",
										"description": "Video Demo Description",
										"video_url": "",
										"cover_url": "http://localhost:3344/raw/image/default/cover.jpg",
										"upload_time": "2019-12-26 14:14:04",
										"update_time": "2019-12-30 21:04:51",
										"author": {
											"uid": 10,
											"username": "aoihosizora",
											"sex": "male",
											"profile": "Demo Profile",
                    						"avatar_url": "http://localhost:3344/raw/image/10/avatar_20191230205334869693.jpg",
											"birth_time": "2019-12-26",
											"authority": "admin"
										}
									}
								]
							}
 						} */
func (v *videoController) QueryVideosByUid(c *gin.Context) {
	uid, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()))
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()))
		return
	}
	page = xconditions.IfThenElse(page < 1, 1, page).(int)

	videos, count := v.videoDao.QueryByUid(uid, page)
	if videos == nil {
		c.JSON(http.StatusNotFound, common.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()))
		return
	}

	c.JSON(http.StatusOK, common.Result{}.Ok().SetPage(count, page, dto.VideoDto{}.FromPos(videos)))
}

// @Router 				/video/{vid} [GET]
// @Summary 			查询视频
// @Description 		查询视频信息
// @Tag					Video
// @Param 				vid path integer true "视频id"
// @Accept 				multipart/form-data
// @ErrorCode			400 request param error
// @ErrorCode			404 video not found
/* @Success 200 		{
							"code": 200,
							"message": "success",
							"data": {
								"vid": 10,
								"title": "Video Title",
								"description": "Video Demo Description",
								"video_url": "",
								"cover_url": "http://localhost:3344/raw/image/default/cover.jpg",
								"upload_time": "2019-12-26 14:14:04",
								"update_time": "2019-12-30 21:04:51",
								"author": {
									"uid": 10,
									"username": "aoihosizora",
									"sex": "male",
									"profile": "Demo Profile",
									"avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg",
									"birth_time": "2019-12-26",
									"authority": "admin"
								}
							}
 						} */
func (v *videoController) QueryVideoByVid(c *gin.Context) {
	vid, err := strconv.Atoi(c.Param("vid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()))
		return
	}

	video := v.videoDao.QueryByVid(vid)
	if video == nil {
		c.JSON(http.StatusNotFound, common.Result{}.Error(http.StatusNotFound).SetMessage(exception.VideoNotFoundError.Error()))
		return
	}

	c.JSON(http.StatusOK, common.Result{}.Ok().SetData(dto.VideoDto{}.FromPo(video)))
}

// @Router 				/video/ [POST] [Auth]
// @Summary 			新建视频
// @Description 		新建用户视频
// @Tag					Video
// @Param 				title formData string true "视频标题" minLength(5) maxLength(100)
// @Param 				description formData string true "视频简介" minLength(0) maxLength(255)
// @Param 				video_url formData string true "视频资源链接"
// @Param 				cover formData file false "视频封面"
// @Accept 				multipart/form-data
// @ErrorCode 			400 request param error
// @ErrorCode 			400 request format error
// @ErrorCode 			400 request body too large
// @ErrorCode 			400 image type not supported
// @ErrorCode 			400 video resource has been used
// @ErrorCode 			500 image save failed
// @ErrorCode 			500 video insert failed
/* @Success 200 		{
							"code": 200,
							"message": "success",
							"data": {
								"vid": 10,
								"title": "Video Title",
								"description": "Video Demo Description",
								"video_url": "",
								"cover_url": "http://localhost:3344/raw/image/default/cover.jpg",
								"upload_time": "2019-12-26 14:14:04",
								"update_time": "2019-12-30 21:04:51",
								"author": {
									"uid": 10,
									"username": "aoihosizora",
									"sex": "male",
									"profile": "Demo Profile",
									"avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg",
									"birth_time": "2019-12-26",
									"authority": "admin"
								}
							}
 						} */
func (v *videoController) InsertVideo(c *gin.Context) {
	authUser := middleware.GetAuthUser(c, v.config)

	videoParam := param.VideoParam{}
	if err := c.ShouldBind(&videoParam); err != nil {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()))
		return
	}
	if !model.FormatCheck.VideoTitle(videoParam.Title) || !model.FormatCheck.VideoDesc(videoParam.Description) {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestFormatError.Error()))
		return
	}
	coverUrl := ""
	coverFile, coverHeader, err := c.Request.FormFile("cover")
	if err == nil && coverFile != nil {
		supported, ext := util.ImageUtil.CheckImageExt(coverHeader.Filename)
		if !supported {
			c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.ImageNotSupportedError.Error()))
			return
		}
		filename := fmt.Sprintf("cover_%s.jpg", util.CommonUtil.CurrentTimeUuid())
		savePath := fmt.Sprintf("./usr/image/%d/%s", authUser.Uid, filename)
		if err := util.ImageUtil.SaveAsJpg(coverFile, ext, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, common.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.ImageSaveError.Error()))
			return
		}
		coverUrl = filename
	}

	video := &po.Video{
		Title:       videoParam.Title,
		Description: videoParam.Description,
		CoverUrl:    coverUrl,
		VideoUrl:    videoParam.VideoUrl, // TODO
		UploadTime:  common.JsonDateTime(time.Now()),
		AuthorUid:   authUser.Uid,
		Author:      authUser,
	}
	status := v.videoDao.Insert(video)
	if status == database.DbExisted {
		c.JSON(http.StatusInternalServerError, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.VideoExistError.Error()))
		return
	} else if status == database.DbFailed {
		c.JSON(http.StatusInternalServerError, common.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.VideoInsertError.Error()))
		return
	}

	c.JSON(http.StatusOK, common.Result{}.Ok().SetData(dto.VideoDto{}.FromPo(video)))
}

// @Router 				/video/{vid} [POST] [Auth]
// @Summary 			更新视频
// @Description 		更新用户视频信息
// @Tag					Video
// @Param 				vid path string true "更新视频id"
// @Param 				title formData string false "视频标题" minLength(5) maxLength(100)
// @Param 				description formData string false "视频简介" minLength(0) maxLength(255)
// @Param 				cover formData file false "视频封面"
// @Accept 				multipart/form-data
// @ErrorCode 			400 request param error
// @ErrorCode 			400 request format error
// @ErrorCode 			400 request body too large
// @ErrorCode 			400 image type not supported
// @ErrorCode 			404 video not found
// @ErrorCode 			500 image save failed
// @ErrorCode 			500 video update failed
/* @Success 200 		{
							"code": 200,
							"message": "success",
							"data": {
								"vid": 10,
								"title": "Video Title",
								"description": "Video Demo Description",
								"video_url": "",
								"cover_url": "http://localhost:3344/raw/image/10/cover_20191230210451040811.jpg",
								"upload_time": "2019-12-26 14:14:04",
								"update_time": "2019-12-30 21:04:51",
								"author": {
									"uid": 10,
									"username": "aoihosizora",
									"sex": "male",
									"profile": "Demo Profile",
									"avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg",
									"birth_time": "2019-12-26",
									"authority": "admin"
								}
							}
 						} */
func (v *videoController) UpdateVideo(c *gin.Context) {
	authUser := middleware.GetAuthUser(c, v.config)

	vid, err := strconv.Atoi(c.Param("vid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()))
		return
	}
	videoParam := param.VideoParam{}
	if err := c.ShouldBind(&videoParam); err != nil {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()))
		return
	}
	if !model.FormatCheck.VideoTitle(videoParam.Title) || !model.FormatCheck.VideoDesc(videoParam.Description) {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestFormatError.Error()))
		return
	}
	coverUrl := ""
	coverFile, coverHeader, err := c.Request.FormFile("cover")
	if err == nil && coverFile != nil {
		supported, ext := util.ImageUtil.CheckImageExt(coverHeader.Filename)
		if !supported {
			c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.ImageNotSupportedError.Error()))
			return
		}
		filename := fmt.Sprintf("cover_%s.jpg", util.CommonUtil.CurrentTimeUuid())
		savePath := fmt.Sprintf("./usr/image/%d/%s", authUser.Uid, filename)
		if err := util.ImageUtil.SaveAsJpg(coverFile, ext, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, common.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.ImageSaveError.Error()))
			return
		}
		coverUrl = filename
	}

	video := v.videoDao.QueryByVid(vid)
	if video == nil {
		c.JSON(http.StatusNotFound, common.Result{}.Error(http.StatusNotFound).SetMessage(exception.VideoNotFoundError.Error()))
		return
	}
	video.Title = videoParam.Title
	video.Description = videoParam.Description
	video.VideoUrl = videoParam.VideoUrl
	video.CoverUrl = xconditions.IfThenElse(coverUrl == "", video.CoverUrl, coverUrl).(string)

	status := v.videoDao.Update(video, authUser.Uid)
	if status == database.DbNotFound {
		c.JSON(http.StatusNotFound, common.Result{}.Error(http.StatusNotFound).SetMessage(exception.VideoNotFoundError.Error()))
		return
	} else if status == database.DbFailed {
		c.JSON(http.StatusInternalServerError, common.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.VideoUpdateError.Error()))
		return
	}

	c.JSON(http.StatusOK, common.Result{}.Ok().SetData(dto.VideoDto{}.FromPo(video)))
}

// @Router 				/video/{vid} [DELETE] [Auth]
// @Summary 			删除视频
// @Description 		删除用户视频
// @Tag					Video
// @Param 				vid path string true "删除视频id"
// @Accept 				multipart/form-data
// @ErrorCode			400 request param error
// @ErrorCode			404 video not found
// @ErrorCode			500 video delete failed
/* @Success 200 		{
							"code": 200,
							"message": "success"
 						} */
func (v *videoController) DeleteVideo(c *gin.Context) {
	authUser := middleware.GetAuthUser(c, v.config)

	vid, err := strconv.Atoi(c.Param("vid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()))
		return
	}

	status := v.videoDao.Delete(vid, authUser.Uid)
	if status == database.DbNotFound {
		c.JSON(http.StatusNotFound, common.Result{}.Error(http.StatusNotFound).SetMessage(exception.VideoNotFoundError.Error()))
		return
	} else if status == database.DbFailed {
		c.JSON(http.StatusInternalServerError, common.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.VideoDeleteError.Error()))
		return
	}

	c.JSON(http.StatusOK, common.Result{}.Ok())
}
