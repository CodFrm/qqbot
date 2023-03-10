package taoke

import (
	"context"
	"strconv"
	"time"

	"github.com/CodFrm/qqbot/utils/taobaoopen"
	"github.com/codfrm/cago/database/redis"
	"github.com/codfrm/cago/pkg/logger"
	"go.uber.org/zap"
)

type repo struct {
}

func (r *repo) IsEnableForward() (bool, error) {
	result, err := redis.Default().Get(context.Background(), "alimama:group:forward:enable").Result()
	if err != nil {
		logger.Default().Error("redis get alimama:group:forward:enable error", zap.Error(err))
		return false, err
	}
	return result == "1", nil
}

func (r *repo) IsOriginGroup(group int64) (bool, error) {
	result, err := redis.Default().HGet(context.Background(),
		"alimama:forward:group", strconv.FormatInt(group, 10)).Result()
	if err != nil {
		logger.Default().Error("redis get alimama:forward:group error", zap.Error(err))
		return false, err
	}
	return result == "1", nil
}

func (r *repo) ForwardGroupList() ([]int64, error) {
	result, err := redis.Default().SMembers(context.Background(), "alimama:group:list").Result()
	if err != nil {
		logger.Default().Error("redis get alimama:group:list error", zap.Error(err))
		return nil, err
	}
	var list []int64
	for _, v := range result {
		group, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			logger.Default().Error("redis get alimama:group:list error", zap.Error(err))
			continue
		}
		list = append(list, group)
	}
	return list, nil
}

func (r *repo) IsTklSend(tkl *taobaoopen.ConverseTkl) bool {
	return !redis.Default().SetNX(context.Background(), "alimama:tkl:is:send:"+tkl.Content[0].TaoID, "1", time.Second*300).Val()
}

func (r *repo) EnableForward() error {
	return redis.Default().Set(context.Background(), "alimama:group:forward:enable", "1", 0).Err()
}

func (r *repo) DisableForward() error {
	return redis.Default().Set(context.Background(), "alimama:group:forward:enable", "0", 0).Err()
}

func (r *repo) AddForwardGroup(group string) error {
	return redis.Default().SAdd(context.Background(), "alimama:group:list", group).Err()
}

func (r *repo) RemoveForwardGroup(group string) error {
	return redis.Default().SRem(context.Background(), "alimama:group:list", group).Err()
}

func (r *repo) AddSourceGroup(group string) error {
	return redis.Default().HSet(context.Background(), "alimama:forward:group", group, "1").Err()
}

func (r *repo) RemoveSourceGroup(group string) error {
	return redis.Default().HDel(context.Background(), "alimama:forward:group", group).Err()
}

func (r *repo) SourceGroupList() ([]int64, error) {
	result, err := redis.Default().HGetAll(context.Background(), "alimama:forward:group").Result()
	if err != nil {
		logger.Default().Error("redis get alimama:forward:group error", zap.Error(err))
		return nil, err
	}
	var list []int64
	for k, v := range result {
		if v == "1" {
			group, err := strconv.ParseInt(k, 10, 64)
			if err != nil {
				logger.Default().Error("redis get alimama:forward:group error", zap.Error(err))
				continue
			}
			list = append(list, group)
		}
	}
	return list, nil
}

func (r *repo) SubscribeQQ() (map[int64]int64, error) {
	result, err := redis.Default().HGetAll(context.Background(), "alimama:subscribe:list").Result()
	if err != nil {
		logger.Default().Error("redis get alimama:subscribe:list error", zap.Error(err))
		return nil, err
	}
	var list = make(map[int64]int64)
	for k, v := range result {
		qq, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			logger.Default().Error("redis get alimama:subscribe:list error", zap.String("k", k), zap.Error(err))
			continue
		}
		group, _ := strconv.ParseInt(v, 10, 64)
		list[qq] = group
	}
	return list, nil
}

func (r *repo) SubscribeTopic(qq int64) ([]string, error) {
	result, err := redis.Default().SMembers(context.Background(), "alimama:subscribe:topic:"+strconv.FormatInt(qq, 10)).Result()
	if err != nil {
		logger.Default().Error("redis get alimama:subscribe:topic error", zap.Error(err))
		return nil, err
	}
	return result, nil
}

func (r *repo) Subscribe(id int64, group int64, i string) error {
	if err := redis.Default().HSet(context.Background(), "alimama:subscribe:list", strconv.FormatInt(id, 10), strconv.FormatInt(group, 10)).Err(); err != nil {
		return err
	}
	return redis.Default().SAdd(context.Background(), "alimama:subscribe:topic:"+strconv.FormatInt(id, 10), i).Err()
}

func (r *repo) UnSubscribe(id int64, topic string) error {
	if err := redis.Default().SRem(context.Background(), "alimama:subscribe:topic:"+strconv.FormatInt(id, 10), topic).Err(); err != nil {
		return err
	}
	// 判断是否还有订阅
	result, err := redis.Default().SMembers(context.Background(), "alimama:subscribe:topic:"+strconv.FormatInt(id, 10)).Result()
	if err != nil {
		logger.Default().Error("redis get alimama:subscribe:topic error", zap.Error(err))
		return err
	}
	if len(result) == 0 {
		return redis.Default().HDel(context.Background(), "alimama:subscribe:list", strconv.FormatInt(id, 10)).Err()
	}
	return nil
}
