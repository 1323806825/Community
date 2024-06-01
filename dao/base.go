package dao

import (
	"blog/dto"
	"blog/global"
	"blog/global/constant"
	"blog/model"
	"blog/utils"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type BaseDao struct {
	Orm *gorm.DB
}

func NewBaseDao() BaseDao {
	return BaseDao{
		Orm: global.DB,
	}
}

// SaveFile 保存文件
func (m *UserDao) SaveFile(info dto.InnerFileInfo, uid uint) (model.UploadFile, error) {
	var iFile model.UploadFile
	info.ConveyFromInnerFile(&iFile, uid)
	if err := m.Orm.Model(&model.UploadFile{}).Save(&iFile).Error; err != nil {
		return iFile, err
	}
	return iFile, nil
}

func (m *BaseDao) Paginate(p dto.Paginate) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset((p.PageNum - 1) * p.PageSize).Limit(p.PageSize)
	}
}

func (m *BaseDao) GenerateTempId(id int) (string, error) {
	if id == 0 {
		return "", errors.New("id is 0")
	}

	strMap, err := global.RedisClient.GetKeysAndValue("*:md-t-" + strconv.Itoa(id) + "-*")
	if err != nil {
		return "", err
	}

	var tempLink string
	if len(strMap) > 0 {
		for k := range strMap {
			tempLink = strings.Split(k, "-")[3]
			break
		}
		return tempLink, nil
	} else {
		tempLink = utils.GenerateRandomString(8)
		redisKey := utils.GeneralRedisKey(constant.QuestionMediaDetailKey, "{id}", strconv.Itoa(id), viper.GetString("custom.prefix"))
		redisKey = utils.GeneralRedisKey(redisKey, "{str}", tempLink, nil)

		var picture model.UploadFile
		if err := global.DB.Where("id = ?", id).First(&picture).Error; err != nil {
			return "", err
		} else {
			fmt.Println("hit db")
			pictureDTO := dto.InnerFileInfo{
				FileName: picture.FileName,
				FilePath: picture.FilePath,
				FileType: picture.FileType,
				FileSize: picture.FileSize,
			}
			js, _ := json.Marshal(pictureDTO)
			_ = global.RedisClient.Set(redisKey, js)
			return tempLink, nil
		}

	}

}
