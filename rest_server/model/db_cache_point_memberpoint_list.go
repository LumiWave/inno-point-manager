package model

import (
	"strconv"
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
)

// redis member point lock key generate
func MakeMemberPointListLockKey(MUID int64) string {
	//return config.GetInstance().DBPrefix + "-MEMBER-POINT-" + strconv.FormatInt(MUID, 10) + "-lock"
	return strconv.FormatInt(MUID, 10) + "-point-lock"
}

// redis member point key generate
func MakeMemberPointListKey(MUID int64) string {
	return config.GetInstance().DBPrefix + ":MEMBER-POINT:" + strconv.FormatInt(MUID, 10)
}

func (o *DB) GetCacheMemberPointList(key string) (*context.PointInfo, error) {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	pointInfos := new(context.PointInfo)

	err := o.Cache.Get(key, pointInfos)

	return pointInfos, err
}

func (o *DB) SetCacheMemberPointList(key string, pointInfo *context.PointInfo) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	conf := config.GetInstance()
	return o.Cache.Set(key, pointInfo, time.Duration(conf.PManager.CachePointExpiryPeriod*int64(time.Minute)))
}

func (o *DB) DelCacheMemberPointList(key string) error {
	return o.Cache.Del(key)
}
