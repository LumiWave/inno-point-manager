package model

import (
	"errors"
	"strconv"
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
)

// redis point lock key generate
func MakePointLockKey(MUID int64) string {
	return config.GetInstance().DBPrefix + "-POINT-MEMBER-" + strconv.FormatInt(MUID, 10) + "-lock"
}

// redis point key generate
func MakePointKey(MUID int64) string {
	return config.GetInstance().DBPrefix + ":POINT-MEMBER:" + strconv.FormatInt(MUID, 10)
}

func AutoLock(key string) (func() error, error) {
	unLock, err := GetDB().Cache.AutoLock(key, nil)
	if err != nil {
		log.Errorf("Result_RedisError_Lock_fail : %v", err)
		return nil, errors.New(resultcode.ResultCodeText[resultcode.Result_RedisError_Lock_fail])
	}

	if unLock == nil {
		log.Errorf("Result_RedisError_Lock_fail : unLock is nil")
		return nil, errors.New(resultcode.ResultCodeText[resultcode.Result_RedisError_Lock_fail])
	}

	return unLock, nil
}

func (o *DB) GetCachePoint(key string) (*context.PointInfo, error) {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	pointInfos := new(context.PointInfo)

	err := o.Cache.Get(key, pointInfos)

	return pointInfos, err
}

func (o *DB) SetCachePoint(key string, pointInfo *context.PointInfo) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	conf := config.GetInstance()
	return o.Cache.Set(key, pointInfo, time.Duration(conf.PManager.CachePointExpiryPeriod*int64(time.Minute)))
}
