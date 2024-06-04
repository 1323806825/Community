package service

import (
	"blog/dao"
	"blog/dto"
	"blog/utils"
)

var subscribeService *SubscribeService

type SubscribeService struct {
	BaseService
	Dao *dao.SubscribeDao
}

func NewSubscribeService() *SubscribeService {
	if subscribeService == nil {
		subscribeService = &SubscribeService{
			Dao: dao.NewSubscribeDao(),
		}
	}
	return subscribeService
}

func (m *SubscribeService) AddSubscribe(iAddSubscribe *dto.AddSubscribeDTO, iUserAuthDTO *dto.UserAuthDTO) error {
	if err := iAddSubscribe.Validate(); err != nil {
		return utils.ParseValidateError(err)
	}
	_, err := m.Dao.GetQuestionDetail(iAddSubscribe.QuestionID, iUserAuthDTO.Uid)
	if err != nil {
		return nil
	}
	err = m.Dao.AddSubscribe(iAddSubscribe, iUserAuthDTO.Uid)
	if err != nil {
		return err
	}
	return m.Dao.AddQuestionSubscribe(iAddSubscribe)

}

func (m *SubscribeService) DeleteSubscribe(iDeleteSubscribe *dto.DeleteSubscribeDTO, iUserAuthDTO *dto.UserAuthDTO) error {
	if err := iDeleteSubscribe.Validate(); err != nil {
		return utils.ParseValidateError(err)
	}

	err := m.Dao.DeleteSubscribe(iDeleteSubscribe, iUserAuthDTO.Uid)
	if err != nil {
		return err
	}
	err = m.Dao.DeleteQuestionSubscribe(iDeleteSubscribe)
	if err != nil {
		return err
	}
	return nil
}
