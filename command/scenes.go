package command

import (
	"errors"
	"github.com/CodFrm/iotqq-plugins/db"
	"github.com/CodFrm/iotqq-plugins/model"
	"github.com/go-redis/redis/v7"
	"strconv"
	"time"
)

var scenes *db.Scenes

func ScenesList(keyword string, page int) ([]string, error) {
	if page <= 0 {
		page = 1
	}
	list, err := scenes.GetList(keyword, page)
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0)
	for _, v := range list {
		ret = append(ret, v.Name)
	}
	return ret, nil
}

func QueryGroupScenes(group int) ([]string, error) {
	return db.Redis.ZRange("group:scenes:"+strconv.FormatInt(int64(group), 10), 0, -1).Result()
}

func QueryScenesTag(name string) (map[string]string, error) {
	ret := make(map[string]string)
	if err := db.GetOrSet("scenes:tag:"+name, &ret, func() (interface{}, error) {
		s, err := scenes.FindScenesByName(name)
		if err != nil {
			return nil, err
		}
		if s == nil {
			return nil, errors.New("场景不存在")
		}
		list, err := scenes.QueryScenesTagByScenesId(s.ID)
		if err != nil {
			return nil, err
		}
		for _, v := range list {
			ret[v.Key] = v.Value
		}
		return ret, nil
	}, db.WithTTL(time.Hour)); err != nil {
		return nil, err
	}
	return ret, nil
}

func QueryMap(qqgroup int, name string) (string, error) {
	list, err := QueryGroupScenes(qqgroup)
	if err != nil {
		return "", err
	}
	for _, v := range list {
		m, err := QueryScenesTag(v)
		if err != nil {
			return "", err
		}
		if v, ok := m[name]; ok {
			return v, nil
		}
	}
	return "", nil
}

func AddScenes(group int, name []string) error {
	list := make([]*redis.Z, 0)
	for _, v := range name {
		if s, err := scenes.FindScenesByName(v); err != nil {
			return err
		} else if s == nil {
			return errors.New(v + "不存在")
		}
		list = append(list, &redis.Z{
			Score:  0,
			Member: v,
		})
	}
	return db.Redis.ZAdd("group:scenes:"+strconv.FormatInt(int64(group), 10), list...).Err()
}

func RemoveScenes(group int, name []string) error {
	return db.Redis.ZRem("group:scenes:"+strconv.FormatInt(int64(group), 10), name).Err()
}

func IsScenesOk(group int) (bool, error) {
	v := db.Redis.Get("scenes:auth:" + strconv.Itoa(group)).Val()
	return v == "1", nil
}

func CreateScenes(scenes_name string, stat int) error {
	m, err := scenes.FindScenesByName(scenes_name)
	if err != nil {
		return err
	}
	if m != nil {
		return errors.New("场景存在")
	}
	m = &model.Scenes{
		Name: scenes_name,
		Stat: int64(stat),
	}
	if err := db.Db.Save(m).Error; err != nil {
		return err
	}
	return nil
}

func AddScenesMap(scenes_name string, key string, value string) error {
	m, err := scenes.FindScenesByName(scenes_name)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("场景不存在")
	}
	if m.Stat == 0 {
		return errors.New("场景已被移除")
	}
	t, err := scenes.FindScenesTag(m.ID, key)
	if err != nil {
		return err
	}
	if t != nil {
		t.Value = value
	} else {
		t = &model.ScenesTag{
			ScenesId: m.ID,
			Key:      key,
			Value:    value,
		}
	}
	if err := db.Db.Save(t).Error; err != nil {
		return err
	}
	db.Redis.Del("scenes:tag:" + scenes_name)
	return nil
}
