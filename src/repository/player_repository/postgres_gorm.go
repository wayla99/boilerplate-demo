package player_repository

import (
	"boilerplate-demo/src/entity/player"
	"boilerplate-demo/src/use_case/repository"
	"context"
	"fmt"
	"strconv"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type gormPostgres struct {
	db     *gorm.DB
	logger *zap.Logger
}

func (g *gormPostgres) Name() string {
	return "gorm_player"
}

func (g *gormPostgres) initSchema() error {
	if err := g.db.AutoMigrate(&postgresGormPlayerModel{}); err != nil {
		return fmt.Errorf("failed to auto migrate: %v", err)
	}

	if err := g.db.AutoMigrate(&postgresGormPointModel{}); err != nil {
		return fmt.Errorf("failed to auto migrate: %v", err)
	}

	return nil
}

func NewGormPostgres(db *gorm.DB, logger *zap.Logger) repository.PlayerRepository {
	g := &gormPostgres{db: db}

	if err := g.initSchema(); err != nil {
		logger.Error("Error initializing schema", zap.Error(err))
	}

	return g
}

func (g *gormPostgres) HealthCheck(ctx context.Context) error {
	if err := g.db.WithContext(ctx).Exec("SELECT 1").Error; err != nil {
		return fmt.Errorf("database health check failed: %v", err)
	}
	return nil
}

func (g *gormPostgres) CreatePlayer(ctx context.Context, player player.Player) (string, error) {
	p, err := playerToGormPlayer(player)
	if err != nil {
		return "", fmt.Errorf("failed to convert player to gorm player: %w", err)
	}
	if err := g.db.WithContext(ctx).Create(&p).Error; err != nil {
		return "", fmt.Errorf("failed to create player: %v", err)
	}

	return strconv.FormatUint(uint64(p.PlayerId), 10), nil
}

func (g *gormPostgres) GetPlayer(ctx context.Context, pId int64) (player.Player, error) {
	var p postgresGormPlayerModel
	if err := g.db.Model(&postgresGormPlayerModel{}).WithContext(ctx).Where("player_id = ? ", pId).First(&p).Error; err != nil {
		return player.Player{}, err
	}

	return p.toEntity(), nil
}
