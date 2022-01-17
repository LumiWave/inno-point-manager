package model

import (
	"strconv"
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/context"
)

// redis member point lock key generate
func MakeDailyAppPointLockKey(appId, pointId int64) string {
	return config.GetInstance().DBPrefix + "-DAILY-APP-POINT-" + strconv.FormatInt(appId, 10) + "-" + strconv.FormatInt(pointId, 10) + "-lock"
}

// redis member point key generate
func MakeDailyAppPointKey(appId, pointId int64) string {
	return config.GetInstance().DBPrefix + ":DAILY-APP-POINT:" + strconv.FormatInt(appId, 10) + "-" + strconv.FormatInt(pointId, 10)
}

func (o *DB) GetCacheDailyAppPoint(key string) (*context.DailyAppPoint, error) {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	dailyAppPoint := new(context.DailyAppPoint)

	err := o.Cache.Get(key, dailyAppPoint)

	return dailyAppPoint, err
}

func (o *DB) SetCacheDailyAppPoint(key string, dailyAppPoint *context.DailyAppPoint) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	conf := config.GetInstance()
	return o.Cache.Set(key, dailyAppPoint, time.Duration(conf.PManager.CachePointExpiryPeriod*int64(time.Minute)))
}
