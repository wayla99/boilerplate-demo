package player_repository

import (
	"boilerplate-demo/src/entity/player"
	"strconv"

	"gorm.io/gorm"
)

type postgresGormPlayerModel struct {
	PlayerId uint   `gorm:"column:player_id;primaryKey;autoIncrement"`
	Name     string `gorm:"column:name"`
	Address  string `gorm:"column:address"`
}

type postgresGormPointModel struct {
	gorm.Model
	Point string
}

func (postgresGormPlayerModel) TableName() string {
	return "players"
}

func (postgresGormPointModel) TableName() string {
	return "points"
}

func (p postgresGormPlayerModel) toEntity() player.Player {
	return player.Player{
		ID:      strconv.FormatUint(uint64(p.PlayerId), 10),
		Name:    p.Name,
		Address: p.Address,
	}
}
