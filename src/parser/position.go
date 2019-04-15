package parser

import (
	"fmt"
	"github.com/siddontang/go-mysql/mysql"
	"go-binlog-replication/src/constants"
	"go-binlog-replication/src/helpers"
	"go-binlog-replication/src/models/system"
	"strconv"
)

var curPosition mysql.Position

func makeHash(dbName string, table string) string {
	return dbName + "." + table
}

func getMinPosition(position mysql.Position) mysql.Position {
	// build current position
	if curPosition.Pos == 0 {
		dbName := helpers.GetCredentials(constants.DBSlave).(helpers.CredentialsDB).DBname

		// get all saved positions for operated tables and fin with min pos
		// WARNING! if it is first start for table, replicate it from min pos of another tables
		for _, table := range helpers.GetTables() {
			hash := makeHash(dbName, table)
			pos, name := helpers.MakeTablePosKey(hash)

			// TODO handle error
			tablePosition, _ := strconv.ParseUint(system.GetValue(pos), 10, 32)
			tableLogFile := system.GetValue(name)

			// TODO check by files, then check by position
			if uint32(tablePosition) < curPosition.Pos || curPosition.Pos == 0 {
				position.Pos = uint32(tablePosition)
				position.Name = tableLogFile
			}
		}
		curPosition = position
	}

	return curPosition
}

// set position for table
func SetPosition(table string, pos mysql.Position) {
	dbName := helpers.GetCredentials(constants.DBSlave).(helpers.CredentialsDB).DBname
	hash := makeHash(dbName, table)

	posKey, nameKey := helpers.MakeTablePosKey(hash)

	system.SetValue(posKey, fmt.Sprint(pos.Pos))
	system.SetValue(nameKey, pos.Name)
}