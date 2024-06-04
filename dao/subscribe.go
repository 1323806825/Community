package dao

import (
	"blog/dto"
	"blog/model"
	"errors"
)

var subscribeDao *SubscribeDao

type SubscribeDao struct {
	BaseDao
}

func NewSubscribeDao() *SubscribeDao {
	if subscribeDao == nil {
		subscribeDao = &SubscribeDao{
			NewBaseDao(),
		}
	}
	return subscribeDao
}

func (m *SubscribeDao) AddSubscribe(stSubscribe *dto.AddSubscribeDTO, uid uint) error {
	var iSubscribe model.Subscribe
	stSubscribe.Convey2Model(&iSubscribe, &uid)
	var count int64
	err := m.Orm.Model(&model.Subscribe{}).Where("question_id = ? and member_id = ?", iSubscribe.QuestionID, iSubscribe.MemberID).Count(&count).Error
	if count > 0 {
		return errors.New("已订阅该问题")
	}
	err = m.Orm.Create(&iSubscribe).Error
	if err != nil {
		return err
	}
	return nil
}

func (m *SubscribeDao) DeleteSubscribe(stSubscribe *dto.DeleteSubscribeDTO, uid uint) error {
	var iSubscribe model.Subscribe
	stSubscribe.Convey2Model(&iSubscribe, &uid)
	var count int64
	m.Orm.Model(&model.Subscribe{}).Where("question_id = ? and member_id = ?", iSubscribe.QuestionID, iSubscribe.MemberID).Count(&count)
	if count < 0 {
		return errors.New("未找到该表")
	}
	err := m.Orm.Model(&model.Subscribe{}).Where("question_id = ?", iSubscribe.QuestionID).Delete(&iSubscribe).Error
	if err != nil {
		return err
	}
	return nil
}

func (m *SubscribeDao) GetQuestionDetail(QuestionID uint, uid uint) (model.Question, error) {
	var detail model.Question
	err := m.Orm.Model(&model.Question{}).Where("id = ?", QuestionID).
		First(&detail).Error
	if err != nil {
		return model.Question{}, err
	}
	if !m.CheckEditQuestionPermission(uid, *detail.OwnerID) {
		return model.Question{}, errors.New("no permission")
	}
	return detail, nil
}

func (m *SubscribeDao) CheckEditQuestionPermission(uid uint, detailId uint) bool {
	return uid == detailId
}

func (m *SubscribeDao) AddQuestionSubscribe(iAddQuestionSubscribe *dto.AddSubscribeDTO) error {
	var Question model.Question
	iAddQuestionSubscribe.ConveyToModel(&Question, &iAddQuestionSubscribe.QuestionID)

	err := m.Orm.Model(&model.Question{}).Where("id = ?", iAddQuestionSubscribe.QuestionID).Update("subscribe_count", Question.SubscribeCount+1).Error
	if err != nil {
		return err
	}
	return nil
}

func (m *SubscribeDao) DeleteQuestionSubscribe(iDeleteQuestionSubscribe *dto.DeleteSubscribeDTO) error {
	var Question model.Question
	iDeleteQuestionSubscribe.ConveyToModel(&Question, &iDeleteQuestionSubscribe.QuestionID)

	err := m.Orm.Model(&model.Question{}).Where("id = ?", iDeleteQuestionSubscribe.QuestionID).Update("subscribe_count", Question.SubscribeCount+1).Error
	return err
}
