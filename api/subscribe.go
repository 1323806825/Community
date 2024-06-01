package api

import (
	"blog/dto"
	"blog/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SubscribeApi struct {
	BaseApi
	Service *service.SubscribeService
}

func NewSubscribeApi() SubscribeApi {
	return SubscribeApi{
		NewBaseApi(),
		service.NewSubscribeService(),
	}
}

func (m *SubscribeApi) AddSubscribe(c *gin.Context) {
	var iAddSubscribeDTO dto.AddSubscribeDTO
	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iAddSubscribeDTO,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	var iUserAuth dto.UserAuthDTO
	iUserAuth.PutAuth(c)

	err := m.Service.AddSubscribe(&iAddSubscribeDTO, &iUserAuth)

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

func (m *SubscribeApi) DeleteSubscribe(c *gin.Context) {
	var iDeleteSubscribeDTO dto.DeleteSubscribeDTO
	if err := m.BuildRequest(BuildRequestOption{
		Ctx:     c,
		DTO:     &iDeleteSubscribeDTO,
		BindAll: true,
	}).GetError(); err != nil {
		return
	}

	var iUserAuth dto.UserAuthDTO
	iUserAuth.PutAuth(c)

	err := m.Service.DeleteSubscribe(&iDeleteSubscribeDTO, &iUserAuth)
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
