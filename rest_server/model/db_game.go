package model

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/LumiWave/baseutil/log"
	orginMssql "github.com/denisenkom/go-mssqldb"
)

const (
	USPG_Mod_Users_ChipQuantity = "[ssrm].[USPG_Mod_Users_ChipQuantity]"
	USPG_Get_Users_By_InnoUID   = "[ssrm].[USPG_Get_Users_By_InnoUID]"
)

func (o *DB) USPG_Mod_Users_ChipQuantity(innoUID string, adjChipQuantity int64, eventID int64) error {
	proc := USPG_Mod_Users_ChipQuantity

	var returnValue orginMssql.ReturnStatus
	rows, err := o.MssqlGameAll.QueryContext(context.Background(), proc,
		sql.Named("InnoUID", innoUID),
		sql.Named("AdjChipQuantity", adjChipQuantity),
		sql.Named("EventID", eventID),
		&returnValue)

	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		log.Errorf("%s QueryContext error : %v", proc, err)
		return err
	}

	if returnValue != 1 {
		log.Errorf("%s returnvalue error : %v", proc, returnValue)
		return errors.New(proc + " returnvalue error " + strconv.Itoa(int(returnValue)))
	}

	return err
}

func (o *DB) USPG_Get_Users_By_InnoUID(innoUID string) (bool, int64, bool, error) {
	proc := USPG_Get_Users_By_InnoUID

	chipQuantity := int64(0)
	isBlocked := false
	var returnValue orginMssql.ReturnStatus
	rows, err := o.MssqlGameRead.QueryContext(context.Background(), proc,
		sql.Named("InnoUID", innoUID),
		&returnValue)

	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		log.Errorf("%s QueryContext error : %v", proc, err)
		return isBlocked, chipQuantity, false, err
	}

	bfind := false
	for rows.Next() {
		if err := rows.Scan(&isBlocked, &chipQuantity); err != nil {
			log.Errorf("USPG_Get_Users_By_InnoUID Scan error : %v", err)
			return isBlocked, chipQuantity, false, err
		}
		bfind = true
	}

	if returnValue != 1 {
		log.Errorf("%s returnvalue error : %v", proc, returnValue)
		return isBlocked, chipQuantity, false, errors.New(proc + " returnvalue error " + strconv.Itoa(int(returnValue)))
	}

	return isBlocked, chipQuantity, bfind, err
}
