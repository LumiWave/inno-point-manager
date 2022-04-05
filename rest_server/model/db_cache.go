package model

import (
	"errors"
	"time"

	"github.com/ONBUFF-IP-TOKEN/basedb"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-point-manager/rest_server/controllers/resultcode"
)

func AutoLock(key string) (func() error, error) {
	opts := new(basedb.LockOptions)
	opts.LockTimeout = 5 * time.Second
	opts.WaitTimeout = 10 * time.Second
	opts.WaitRetry = 10 * time.Millisecond
	unLock, err := GetDB().Cache.AutoLock(key, opts)

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
