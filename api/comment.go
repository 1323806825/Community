package api

import (
	"blog/dto"
	"blog/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CommentApi struct {
	BaseApi
	Service *service.CommentService
}

func NewCommentApi() CommentApi {
	return CommentApi{
		NewBaseApi(),
		service.NewCommentService(),
	}
}

func (m *CommentApi) GetCommentList(c *gin.Context) {
	var iCommentList dto.CommentListDTO
	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iCommentList,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	var nTotal int64
	list, nTotal, err := m.Service.GetContentList(&iCommentList)
	if err != nil {
		m.SetError(err)
		m.Fail(ResponseJson{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	m.OK(ResponseJson{
		Code:    0,
		Message: "success",
		Data:    list,
		Total:   nTotal,
	})
}

func (m *CommentApi) CommentPostQuestion(c *gin.Context) {
	var iPostComment dto.PostComment
	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iPostComment,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	var iUserAuth dto.UserAuthDTO
	iUserAuth.PutAuth(c)

	comment, err := m.Service.PostComment(&iPostComment, &iUserAuth)
	if err != nil {
		m.Fail(ResponseJson{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	type ReData struct {
		QuestionID uint   `json:"question_id"`
		Content    string `json:"content"`
	}

	m.OK(ResponseJson{
		Code:    0,
		Message: "post comment success",
		Data: ReData{
			QuestionID: comment.QuestionID,
			Content:    comment.Content,
		},
	})

}

func (m *CommentApi) CommentPostComment(c *gin.Context) {
	var iCommentPostComment dto.CommentPostComment
	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iCommentPostComment,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	var iUserAuth dto.UserAuthDTO
	iUserAuth.PutAuth(c)

	data, err := m.Service.CommentPostComment(&iCommentPostComment, &iUserAuth)
	if err != nil {
		m.Fail(ResponseJson{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	type ReData struct {
		Content         string `json:"content"`
		ParentCommentID uint   `json:"parent_comment_id"`
	}
	m.OK(ResponseJson{
		Code:    0,
		Message: "post comment success",
		Data: ReData{
			Content:         data.Content,
			ParentCommentID: data.ParentCommentID,
		},
	})
}

func (m *CommentApi) GetMyComment(c *gin.Context) {
	var iCommentList dto.CommentListDTO
	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iCommentList,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	var iUserAuth dto.UserAuthDTO
	iUserAuth.PutAuth(c)

	list, nTotal, err := m.Service.GetMyComment(&iCommentList, &iUserAuth)
	if err != nil {
		m.Fail(ResponseJson{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Code:    0,
		Message: "get my comment success",
		Data:    list,
		Total:   nTotal,
	})
}

func (m *CommentApi) DeleteComment(c *gin.Context) {
	var iDeleteComment dto.DeleteCommentDTO
	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iDeleteComment,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	var iUserAuth dto.UserAuthDTO
	iUserAuth.PutAuth(c)

	err := m.Service.DeleteComment(&iDeleteComment, &iUserAuth)
	if err != nil {
		m.Fail(ResponseJson{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	m.OK(ResponseJson{
		Code:    0,
		Message: "delete success",
	})
}

func (m *CommentApi) GetMyCommentLikeCount(c *gin.Context) {
	var iGetMyCommentLikeCount dto.GetMyCommentLikeCount
	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iGetMyCommentLikeCount,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	var iUserAuth dto.UserAuthDTO
	iUserAuth.PutAuth(c)

	nTotal, err := m.Service.GetMyLikeCount(&iGetMyCommentLikeCount, &iUserAuth)
	if err != nil {
		m.SetError(err)
		m.Fail(ResponseJson{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	m.OK(ResponseJson{
		Code:    0,
		Message: "get like count success",
		Total:   nTotal,
	})
}

func (m *CommentApi) UploadCommentPicture(c *gin.Context) {
	var iUploadCommentPicture dto.UploadCommentPictureDTO

	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iUploadCommentPicture,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	var iUserAuth dto.UserAuthDTO
	iUserAuth.PutAuth(c)

	err := m.Service.UploadCommentPictureFormValidate(&iUploadCommentPicture, &iUserAuth)
	if err != nil {
		m.Fail(ResponseJson{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	err = m.Service.CheckCommentPictureExist(&iUploadCommentPicture)
	if err != nil {
		m.Fail(ResponseJson{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	uploadPath := "./upload/"

	file, err := m.SaveFile("image", uploadPath, FileUploadLimit{
		MaxSize:   1024 * 1024 * 5,
		MaxCount:  1,
		AllowType: []string{"jpg", "jpeg", "png", "svg"},
	})
	if err != nil {
		m.Fail(ResponseJson{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	tempLink, err := m.Service.UploadCommentPicture(iUploadCommentPicture.CommentID, file[0], &iUserAuth)
	if err != nil {
		m.Fail(ResponseJson{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Code:    0,
		Message: "success",
		Data:    tempLink,
	})
}

func (m *CommentApi) DeleteCommentPicture(c *gin.Context) {
	var iDeleteCommentPictureDTO dto.DeleteCommentPictureDTO
	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iDeleteCommentPictureDTO,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	var iUserAuthDTO dto.UserAuthDTO
	iUserAuthDTO.PutAuth(c)

	err := m.Service.DeleteCommentPicture(&iDeleteCommentPictureDTO, &iUserAuthDTO)

	if err != nil {
		m.SetError(err)
		m.Fail(ResponseJson{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	m.OK(ResponseJson{
		Code:    0,
		Message: "success",
	})
}
