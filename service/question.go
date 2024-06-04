package service

import (
	"blog/dao"
	"blog/dto"
	"blog/model"
	"blog/utils"
)

var questionService *QuestionService

type QuestionService struct {
	BaseService
	Dao *dao.QuestionDao
}

func NewQuestionService() *QuestionService {
	if questionService == nil {
		questionService = &QuestionService{
			Dao: dao.NewQuestionDao(),
		}
	}
	return questionService
}

func (m *QuestionService) AddQuestion(iQuestion *dto.QuestionDTO, iUserAuth *dto.UserAuthDTO) (model.Question, error) {
	if err := iQuestion.Validate(); err != nil {
		return model.Question{}, utils.ParseValidateError(err)
	}

	Question, err := m.Dao.AddQuestion(iQuestion, iUserAuth.Uid)
	if err != nil {
		return model.Question{}, err
	}
	return Question, nil
}

func (m *QuestionService) GetQuestionList(iQuestionList *dto.QuestionListDTO) ([]model.Question, int64, error) {
	if err := iQuestionList.Validate(); err != nil {
		return nil, 0, utils.ParseValidateError(err)
	}
	if iQuestionList.PageNum <= 0 {
		iQuestionList.PageNum = 1
	}
	if iQuestionList.PageSize <= 0 {
		iQuestionList.PageSize = 10
	}
	list, nTotal, err := m.Dao.GetQuestionList(iQuestionList.Paginate)
	if err != nil {
		return nil, 0, err
	}
	return list, nTotal, nil
}

func (m *QuestionService) GetQuestionById(uid uint) (model.Question, error) {
	detail, err := m.Dao.GetQuestionById(uid)
	if err != nil {
		return model.Question{}, err
	}
	//var user model.User
	//user, _ = m.Dao.FindUserById(*detail.OwnerID)
	//detail.OwnerID = &user.ID
	return detail, nil
}

func (m *QuestionService) GetMyQuestion(iQuestionList *dto.QuestionListDTO, iUserAuth *dto.UserAuthDTO) ([]model.Question, int64, error) {
	var list []model.Question
	var nTotal int64
	if iQuestionList.PageNum <= 0 {
		iQuestionList.PageNum = 1
	}
	if iQuestionList.PageSize <= 0 {
		iQuestionList.PageSize = 10
	}
	list, nTotal, err := m.Dao.GetMyQuestion(iQuestionList.Paginate, iUserAuth.Uid)
	if err != nil {
		return nil, 0, err
	}
	return list, nTotal, nil
}

func (m *QuestionService) UpdateTitle(iUpdateTitle *dto.UpdateTitleDTO, iUserAuth *dto.UserAuthDTO) error {
	if err := iUpdateTitle.Validate(); err != nil {
		return utils.ParseValidateError(err)
	}
	_, err := m.Dao.GetQuestionDetail(iUpdateTitle.QuestionID, iUserAuth.Uid)
	if err != nil {
		return err
	}
	return m.Dao.UpdateTitle(iUpdateTitle.Title, iUpdateTitle.QuestionID)
}

func (m *QuestionService) UpdateContent(iUpdateContent *dto.UpdateContentDTO, iUserAuth *dto.UserAuthDTO) error {
	if err := iUpdateContent.Validate(); err != nil {
		return utils.ParseValidateError(err)
	}
	_, err := m.Dao.GetQuestionDetail(iUpdateContent.QuestionID, iUserAuth.Uid)
	if err != nil {
		return err
	}
	return m.Dao.UpdateContent(iUpdateContent.Content, iUpdateContent.QuestionID)
}

func (m *QuestionService) CheckQuestionPictureFormValidate(iQuestionPictureDTO *dto.UploadQuestionPictureDTO, uid uint) error {
	if err := iQuestionPictureDTO.Validate(); err != nil {
		return utils.ParseValidateError(err)
	}
	_, err := m.Dao.GetQuestionDetail(iQuestionPictureDTO.QuestionID, uid)
	if err != nil {
		return err
	}
	return nil
}

func (m *QuestionService) CheckQuestionPictureExist(iQuestionPictureDTO *dto.UploadQuestionPictureDTO) error {
	return m.Dao.CheckQuestionPictureExist(iQuestionPictureDTO.QuestionID)
}

func (m *QuestionService) UploadQuestionPicture(questionID uint, innerFile dto.InnerFileInfo, iUserAuth *dto.UserAuthDTO) (string, error) {
	file, err := m.Dao.SaveFile(innerFile, iUserAuth.Uid)
	if err != nil {
		return "", nil
	}
	tempLink, err := m.Dao.GenerateTempId(int(file.ID))
	if err != nil {
		return "", err
	}
	err = m.Dao.SetQuestionPicture(file.ID, questionID)
	if err != nil {
		return "", err
	}
	return tempLink, nil
}

func (m *QuestionService) DeleteQuestion(iDeleteQuestion *dto.DeleteQuestionDTO, iUserAuth *dto.UserAuthDTO) error {
	if err := iDeleteQuestion.Validate(); err != nil {
		return utils.ParseValidateError(err)
	}
	question, err := m.Dao.GetQuestionDetail(iDeleteQuestion.QuestionID, iUserAuth.Uid)
	if err != nil {
		return err
	}
	return m.Dao.DeleteQuestion(iDeleteQuestion.QuestionID, question)
}

func (m *QuestionService) DeleteQuestionPicture(iDeleteQuestionPicture *dto.DeleteQuestionPictureDTO, iUserAuthDTO *dto.UserAuthDTO) error {
	if err := iDeleteQuestionPicture.Validate(); err != nil {
		return utils.ParseValidateError(err)
	}
	_, err := m.Dao.GetQuestionDetail(iDeleteQuestionPicture.QuestionID, iUserAuthDTO.Uid)
	if err != nil {
		return err
	}
	return m.Dao.DeleteQuestionPicture(iDeleteQuestionPicture.QuestionID)
}

func (m *QuestionService) AddQuestionSubscribe(iAddSubscribeDTO *dto.AddQuestionSubscribeDTO, iUserAuthDTO *dto.UserAuthDTO) error {
	if err := iAddSubscribeDTO.Validate(); err != nil {
		return utils.ParseValidateError(err)
	}

	_, err := m.Dao.GetQuestionDetail(iAddSubscribeDTO.QuestionID, iUserAuthDTO.Uid)
	if err != nil {
		return err
	}

	return m.Dao.AddQuestionSubscribe(iAddSubscribeDTO)

}

func (m *QuestionService) GetMySubscribe(iGetMySubscribeDTO *dto.GetMySubscribeDTO, iUserAuthDTO *dto.UserAuthDTO) (int64, error) {
	if err := iGetMySubscribeDTO.Validate(); err != nil {
		return 0, utils.ParseValidateError(err)
	}

	_, err := m.Dao.GetQuestionDetail(iGetMySubscribeDTO.QuestionID, iUserAuthDTO.Uid)
	if err != nil {
		return 0, err
	}
	total, err := m.Dao.GetMySubscribe(iGetMySubscribeDTO.QuestionID)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (m *QuestionService) DeleteQuestionSubscribe(iDeleteSubscribeDTO *dto.DeleteQuestionSubscribeDTO, iUserAuthDTO *dto.UserAuthDTO) error {
	if err := iDeleteSubscribeDTO.Validate(); err != nil {
		return utils.ParseValidateError(err)
	}

	_, err := m.Dao.GetQuestionDetail(iDeleteSubscribeDTO.QuestionID, iUserAuthDTO.Uid)
	if err != nil {
		return err
	}

	return m.Dao.DeleteQuestionSubscribe(iDeleteSubscribeDTO)
}
