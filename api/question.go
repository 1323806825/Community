package api

import (
	"blog/dto"
	"blog/model"
	"blog/service"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type QuestionApi struct {
	BaseApi
	Service *service.QuestionService
}

func NewQuestionApi() QuestionApi {
	return QuestionApi{
		NewBaseApi(),
		service.NewQuestionService(),
	}
}

func (m *QuestionApi) AddQuestion(c *gin.Context) {
	var iQuestionDTO dto.QuestionDTO
	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iQuestionDTO,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	var iUserAuth dto.UserAuthDTO
	iUserAuth.PutAuth(c)

	Question, err := m.Service.AddQuestion(&iQuestionDTO, &iUserAuth)
	if err != nil {
		m.SetError(err)
		m.Fail(ResponseJson{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	type reData struct {
		QuestionID uint   `json:"question_id"`
		Title      string `json:"title"`
		Content    string `json:"content"`
	}
	m.OK(ResponseJson{
		Code:    0,
		Message: "post question success",
		Data: reData{
			QuestionID: Question.ID,
			Title:      Question.Title,
			Content:    Question.Content,
		},
	})
}

func (m *QuestionApi) GetQuestionList(c *gin.Context) {
	var iQuestionList dto.QuestionListDTO
	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iQuestionList,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	var nTotal int64
	list, nTotal, err := m.Service.GetQuestionList(&iQuestionList)
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

func (m *QuestionApi) GetQuestionById(c *gin.Context) {
	if err := m.BuildRequest(BuildRequestOption{
		Ctx: c,
	}).GetError(); err != nil {
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		m.OK(ResponseJson{
			Code:    http.StatusBadRequest,
			Message: errors.New("id must be a number").Error(),
		})
		return
	}

	detail, err := m.Service.GetQuestionById(uint(id))
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		m.OK(ResponseJson{
			Data: []model.Question{},
		})
		return
	} else if err != nil {
		m.OK(ResponseJson{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	m.OK(ResponseJson{
		Code:    0,
		Message: "success",
		Data:    detail,
	})
}

func (m *QuestionApi) GetMyQuestion(c *gin.Context) {
	var iUserAuthDTO dto.UserAuthDTO
	iUserAuthDTO.PutAuth(c)

	var iQuestionList dto.QuestionListDTO
	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iQuestionList,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	list, nTotal, err := m.Service.GetMyQuestion(&iQuestionList, &iUserAuthDTO)
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

func (m *QuestionApi) UpdateTitle(c *gin.Context) {
	var iUpdateTitleDTO dto.UpdateTitleDTO
	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iUpdateTitleDTO,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	var iUserAuthDTO dto.UserAuthDTO
	iUserAuthDTO.PutAuth(c)

	err := m.Service.UpdateTitle(&iUpdateTitleDTO, &iUserAuthDTO)
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

func (m *QuestionApi) UpdateContent(c *gin.Context) {
	var iUpdateContentDTO dto.UpdateContentDTO
	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iUpdateContentDTO,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	var iUserAuthDTO dto.UserAuthDTO
	iUserAuthDTO.PutAuth(c)

	err := m.Service.UpdateContent(&iUpdateContentDTO, &iUserAuthDTO)
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

func (m *QuestionApi) UploadPicture(c *gin.Context) {
	var iUploadQuestionPicture dto.UploadQuestionPictureDTO
	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iUploadQuestionPicture,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	var iUserAuth dto.UserAuthDTO
	iUserAuth.PutAuth(c)

	//判断图片格式合法化
	err := m.Service.CheckQuestionPictureFormValidate(&iUploadQuestionPicture, iUserAuth.Uid)
	if err != nil {
		m.OK(ResponseJson{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	//检查图片是否存在
	err = m.Service.CheckQuestionPictureExist(&iUploadQuestionPicture)
	if err != nil {
		m.OK(ResponseJson{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	//运行文件目录下upload文件夹
	uploadPath := "./upload/"

	file, err := m.SaveFile("image", uploadPath, FileUploadLimit{
		MaxCount:  1,
		MaxSize:   1024 * 1024 * 5,
		AllowType: []string{"jpg", "jpeg", "png", "svg"},
	})

	if err != nil {
		m.OK(ResponseJson{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	//上传问题文件
	tempLink, err := m.Service.UploadQuestionPicture(iUploadQuestionPicture.QuestionID, file[0], &iUserAuth)
	if err != nil {
		m.OK(ResponseJson{
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

func (m *QuestionApi) DeleteQuestionPicture(c *gin.Context) {
	var iDeleteQuestionPictureDTO dto.DeleteQuestionPictureDTO
	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iDeleteQuestionPictureDTO,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	var iUserAuthDTO dto.UserAuthDTO
	iUserAuthDTO.PutAuth(c)

	err := m.Service.DeleteQuestionPicture(&iDeleteQuestionPictureDTO, &iUserAuthDTO)
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
		Message: "delete picture success",
	})
}

func (m *QuestionApi) DeleteQuestion(c *gin.Context) {
	var iDeleteQuestion dto.DeleteQuestionDTO
	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iDeleteQuestion,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	var iUserAuth dto.UserAuthDTO
	iUserAuth.PutAuth(c)

	err := m.Service.DeleteQuestion(&iDeleteQuestion, &iUserAuth)
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

func (m *QuestionApi) AddQuestionSubscribe(c *gin.Context) {
	var iAddQuestionSubscribe dto.AddQuestionSubscribeDTO
	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iAddQuestionSubscribe,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	var iUserAuth dto.UserAuthDTO
	iUserAuth.PutAuth(c)

	err := m.Service.AddQuestionSubscribe(&iAddQuestionSubscribe, &iUserAuth)
	if err != nil {
		m.Fail(ResponseJson{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	m.OK(ResponseJson{
		Code:    0,
		Message: "add subscribe success",
	})
}

func (m *QuestionApi) DeleteQuestionSubscribe(c *gin.Context) {
	var iDeleteQuestionSubscribe dto.DeleteQuestionSubscribeDTO
	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iDeleteQuestionSubscribe,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	var iUserAuth dto.UserAuthDTO
	iUserAuth.PutAuth(c)

	err := m.Service.DeleteQuestionSubscribe(&iDeleteQuestionSubscribe, &iUserAuth)
	if err != nil {
		m.Fail(ResponseJson{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	m.OK(ResponseJson{
		Code:    0,
		Message: "delete subscribe success",
	})
}

func (m *QuestionApi) GetMyQuestionSubscribe(c *gin.Context) {
	var iGetMySubscribe dto.GetMySubscribeDTO
	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iGetMySubscribe,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	var iUserAuth dto.UserAuthDTO
	iUserAuth.PutAuth(c)

	total, err := m.Service.GetMySubscribe(&iGetMySubscribe, &iUserAuth)
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
		Total:   total,
	})
}
