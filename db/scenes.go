package db

import (
	"github.com/CodFrm/iotqq-plugins/model"
	"github.com/jinzhu/gorm"
)

type Scenes struct {
}

func NewScenes() *Scenes {
	return &Scenes{}
}

func (s *Scenes) GetList(keyword string, page int) ([]*model.Scenes, error) {
	ret := make([]*model.Scenes, 0)
	db := Db.Model(model.Scenes{}).Order("id desc").Where("stat=1")
	if keyword != "" {
		db = db.Where("name like ?", "%"+keyword+"%")
	}
	if err := db.Limit(20).Offset((page - 1) * 20).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *Scenes) FindScenesByName(name string) (*model.Scenes, error) {
	ret := &model.Scenes{}
	if err := Db.Where("name=?", name).First(ret).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (s *Scenes) QueryScenesTagByScenesId(scenes_id int64) ([]*model.ScenesTag, error) {
	ret := make([]*model.ScenesTag, 0)
	db := Db.Model(model.ScenesTag{}).Where("scenes_id=?", scenes_id)
	if err := db.Scan(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}
