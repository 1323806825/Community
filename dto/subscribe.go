package dto

import (
	"blog/model"
	"github.com/go-playground/validator/v10"
)

type AddSubscribeDTO struct {
	QuestionID uint `json:"question_id" form:"question_id" validate:"required,number"`
}

func (m *AddSubscribeDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}

func (m *AddSubscribeDTO) Convey2Model(subscribe *model.Subscribe, uid *uint) {
	subscribe.MemberID = uid
	subscribe.QuestionID = &m.QuestionID
}

func (m *AddSubscribeDTO) ConveyToModel(AddQuestionSubscribe *model.Question, uid *uint) {
	AddQuestionSubscribe.OwnerID = uid

}

type DeleteSubscribeDTO struct {
	QuestionID uint `json:"question_id" form:"question_id" validate:"required,number"`
}

func (m *DeleteSubscribeDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}

func (m *DeleteSubscribeDTO) Convey2Model(subscribe *model.Subscribe, uid *uint) {
	subscribe.MemberID = uid
	subscribe.QuestionID = &m.QuestionID
}

func (m *DeleteSubscribeDTO) ConveyToModel(DeleteQuestionSubscribe *model.Question, uid *uint) {
	DeleteQuestionSubscribe.OwnerID = uid

}
