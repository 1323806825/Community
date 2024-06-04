package dao

import (
	"blog/dto"
	"blog/model"
	"errors"
	"gorm.io/gorm"
	"os"
)

var questionDao *QuestionDao

type QuestionDao struct {
	BaseDao
}

func NewQuestionDao() *QuestionDao {
	if questionDao == nil {
		questionDao = &QuestionDao{
			NewBaseDao(),
		}
	}
	return questionDao
}

func (m *QuestionDao) AddQuestion(QuestionDTO *dto.QuestionDTO, uid uint) (model.Question, error) {
	var iQuestion model.Question
	QuestionDTO.ConveyToModel(&iQuestion, &uid)
	err := m.Orm.Save(&iQuestion).Error
	if err != nil {
		return model.Question{}, err
	}
	return iQuestion, nil
}

func (m *QuestionDao) Paginate(p dto.Paginate) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset((p.PageNum - 1) * p.PageSize).Limit(p.PageSize)
	}
}

func (m *QuestionDao) GetQuestionList(paginate dto.Paginate) (iQuestionList []model.Question, nTotal int64, err error) {
	//var iQuestionList []model.Question
	err = m.Orm.Model(&model.Question{}).
		Scopes(m.Paginate(paginate)).
		Order("created_at desc").
		Find(&iQuestionList).
		Offset(-1).Limit(-1).
		Count(&nTotal).Error
	if err != nil {
		return nil, 0, err
	}
	return iQuestionList, nTotal, nil
}

func (m *QuestionDao) GetQuestionById(uid uint) (model.Question, error) {
	var detail model.Question
	err := m.Orm.Model(&model.Question{}).Where("id = ?", uid).
		First(&detail).Error
	if err != nil {
		return model.Question{}, err
	}
	return detail, nil
}

func (m *QuestionDao) FindUserById(uid uint) (model.User, error) {
	var iUser model.User
	err := m.Orm.Model(&model.User{}).Where("id = ?", uid).
		First(&iUser).Error
	if err != nil {
		return model.User{}, err
	}
	//var iAuthUser model.AuthUser
	//iAuthUser = iUser.InAuth()
	return iUser, nil
}

func (m *QuestionDao) GetMyQuestion(paginate dto.Paginate, uid uint) ([]model.Question, int64, error) {
	var list []model.Question
	var nTotal int64
	err := m.Orm.Model(&model.Question{}).
		Where("owner_id = ?", uid).
		Scopes(m.Paginate(paginate)).
		Order("created_at desc").
		Find(&list).
		Offset(-1).Limit(-1).
		Count(&nTotal).
		Error
	if err != nil {
		return nil, 0, err
	}
	return list, nTotal, nil
}

func (m *QuestionDao) CheckEditQuestionPermission(uid uint, detailId uint) bool {
	return uid == detailId
}

func (m *QuestionDao) GetQuestionDetail(QuestionID uint, uid uint) (model.Question, error) {
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

func (m *QuestionDao) UpdateTitle(Title string, uid uint) error {
	err := m.Orm.Model(&model.Question{}).Where("id = ?", uid).
		Update("title", Title).Error
	return err
}

func (m *QuestionDao) UpdateContent(Content string, uid uint) error {
	err := m.Orm.Model(&model.Question{}).Where("id = ?", uid).Update("content", Content).Error
	return err
}

func (m *QuestionDao) SaveFile(info dto.InnerFileInfo, uid uint) (model.UploadFile, error) {
	var iFile model.UploadFile
	info.ConveyFromInnerFile(&iFile, uid)
	if err := m.Orm.Model(&model.UploadFile{}).Save(&iFile).Error; err != nil {
		return iFile, err
	}
	return iFile, nil
}

func (m *QuestionDao) SetQuestionPicture(fileID uint, QuestionID uint) error {
	err := m.Orm.Model(&model.Question{}).Where("id = ?", QuestionID).Update("picture_id", fileID).Error
	return err
}

func (m *QuestionDao) CheckQuestionPictureExist(uid uint) error {
	var question model.Question
	err := m.Orm.Model(&model.Question{}).Where("id = ?", uid).First(&question).Error
	if err != nil {
		return nil
	}
	if question.PictureID != nil {
		var pic model.UploadFile
		err = m.Orm.Model(&model.UploadFile{}).Where("id = ?", question.PictureID).Error
		if err != nil {
			return nil
		}
		err = os.Remove(pic.FilePath)
		if err != nil {
			return err
		}
		err = m.Orm.Model(&model.UploadFile{}).Where("id = ?", pic.ID).Delete(&pic).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *QuestionDao) DeleteQuestion(uid uint, question model.Question) error {
	err := m.Orm.Model(&model.Question{}).Where("id = ?", uid).Delete(&question).Error
	if err != nil {
		return err
	}
	return nil
}

func (m *QuestionDao) DeleteQuestionPicture(uid uint) error {
	var iFile model.UploadFile
	err := m.Orm.Model(&model.UploadFile{}).Where("id = ?", uid).Delete(&iFile).Error
	return err
}

func (m *QuestionDao) AddQuestionSubscribe(iAddQuestionSubscribe *dto.AddQuestionSubscribeDTO) error {
	var Question model.Question
	iAddQuestionSubscribe.ConveyToModel(&Question, &iAddQuestionSubscribe.QuestionID)

	err := m.Orm.Model(&model.Question{}).Where("id = ?", iAddQuestionSubscribe.QuestionID).Update("subscribe_count", Question.SubscribeCount+1).Error
	return err
}

func (m *QuestionDao) GetMySubscribe(uid uint) (int64, error) {
	var question model.Question
	err := m.Orm.Model(&model.Question{}).Where("id = ?", uid).First(&question).Error
	if err != nil {
		return 0, err
	}
	return question.SubscribeCount, nil
}

func (m *QuestionDao) DeleteQuestionSubscribe(iDeleteQuestionSubscribe *dto.DeleteQuestionSubscribeDTO) error {
	var Question model.Question
	iDeleteQuestionSubscribe.ConveyToModel(&Question, &iDeleteQuestionSubscribe.QuestionID)

	err := m.Orm.Model(&model.Question{}).Where("id = ?", iDeleteQuestionSubscribe.QuestionID).Update("subscribe_count", Question.SubscribeCount+1).Error
	return err
}
