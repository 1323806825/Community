package api

import (
	"blog/dto"
	"blog/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
)

type BaseApi struct {
	Ctx    *gin.Context
	Errors error
}

func NewBaseApi() BaseApi {
	return BaseApi{}
}

func (m *BaseApi) SetError(err error) {
	m.Errors = err
}

func (m *BaseApi) GetError() error {
	return m.Errors
}

type BuildRequestOption struct {
	Ctx     *gin.Context
	DTO     any
	BindUri bool //用于判断是否从uri取数据
	BindAll bool //所有获取渠道
}

func (m *BaseApi) BuildRequest(option BuildRequestOption) *BaseApi {
	//整合错误
	var errResult error

	//绑定ctx
	m.Ctx = option.Ctx

	//绑定请求数据
	if option.DTO != nil {
		//只有非空才进行绑定
		//判断是否从uri获取数据
		if option.BindUri || option.BindAll {
			errResult = m.Ctx.ShouldBindUri(option.DTO)
		}

		if !option.BindUri || option.BindAll {
			errResult = m.Ctx.ShouldBind(option.DTO)
		}

		//判断是否有误
		if errResult != nil {
			//开始解释这个错误到我们的提示语言
			m.SetError(errResult)
			m.Fail(ResponseJson{
				Code:    400,
				Message: m.GetError().Error(),
			})
		}
	}
	return m
}

func (m *BaseApi) Fail(resp ResponseJson) {
	Fail(m.Ctx, resp)
}

func (m *BaseApi) OK(resp ResponseJson) {
	OK(m.Ctx, resp)
}

func (m *BaseApi) ServerFail(resp ResponseJson) {
	ServerFail(m.Ctx, resp)
}

type FileUploadLimit struct {
	MaxSize   int64
	MaxCount  int
	AllowType []string
}

// SaveFile 处理文件上传，保存进服务器路径
func (m *BaseApi) SaveFile(headerKey string, path string, limit FileUploadLimit) ([]dto.InnerFileInfo, error) {
	form, _ := m.Ctx.MultipartForm()
	files := form.File[headerKey]
	returnList := make([]dto.InnerFileInfo, 0)
	//path是文件保存路径，不包含文件名
	if len(files) == 0 {
		return returnList, errors.New("no file upload")
	}

	//检查文件是否合法
	if len(files) > limit.MaxCount && limit.MaxCount != 0 {
		return returnList, errors.New("too many files")
	}

	//检查文件类型，遍历文件
	haveInvalidType := false
	for _, file := range files {
		//检查文件大小
		if file.Size > limit.MaxSize && limit.MaxSize != 0 {
			return returnList, errors.New("file too large")
		}
		//检查文件类型
		//获取文件后缀
		suffix := strings.Split(file.Filename, ".")[2]

		//检查是否合法类型
		if !checkInvalidType(suffix, limit.AllowType) {
			haveInvalidType = true
			break
		}

	}
	if haveInvalidType {
		return returnList, errors.New("invalid file type")
	}

	//保存文件 save file
	for _, file := range files {
		newFileName := utils.GenerateRandomString(7) + "_" + file.Filename
		dst := path + newFileName

		if err := m.Ctx.SaveUploadedFile(file, dst); err != nil {
			return returnList, err
		} else {
			//保存成功
			returnList = append(returnList, dto.InnerFileInfo{
				FileName: newFileName,
				FilePath: dst,
				FileSize: file.Size,
				FileType: strings.Split(file.Filename, ".")[2],
			})
		}
	}
	return returnList, nil
}

func checkInvalidType(iType string, allowType []string) bool {
	//看iType 是否在allowType中
	for _, v := range allowType {
		if v == iType {
			return true
		}
	}
	return false
}
