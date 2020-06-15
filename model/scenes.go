package model

type Scenes struct {
	ID   int64  `gorm:"column:id" json:"id" form:"id"`
	Name string `gorm:"column:name" json:"name" form:"name"`
	Stat int64  `gorm:"column:stat" json:"stat" form:"stat"`
}

type ScenesTag struct {
	ID       int64  `gorm:"column:id" json:"id" form:"id"`
	ScenesId int64  `gorm:"column:scenes_id" json:"scenes_id" form:"scenes_id"`
	Key      string `gorm:"column:key" json:"key" form:"key"`
	Value    string `gorm:"column:value" json:"value" form:"value"`
}
