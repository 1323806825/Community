package dto

import "blog/model"

type InnerFile struct{}

func (m *InnerFileInfo) ConveyFromInnerFile(upload *model.UploadFile, uid uint) {
	upload.FileSize = m.FileSize
	upload.FileType = m.FileType
	upload.FilePath = m.FilePath
	upload.FileName = m.FileName
	upload.Uid = uid
}
