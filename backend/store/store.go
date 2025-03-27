package store

import (
	"context"
	"fmt"
	"log/slog"
	"rental-property-management-system/backend/config"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	rdb *redis.Client
)

func GetDB() *gorm.DB {
	return db
}

// Redis 缓存层相关操作函数
type RedisNamespace string

const (
	redisSession            RedisNamespace = "session"
	redisExpiredAccessToken RedisNamespace = "expired_access_token"
)

func redisKey(args ...string) string {
	return strings.Join(args, ":")
}

func SetSession(ctx context.Context, key string, jsonValue []byte) error {
	if err := rdb.Set(ctx, redisKey(string(redisSession), key), jsonValue, config.AccessTokenExpiration).Err(); err != nil {
		return errors.Wrapf(err, "failed to set session")
	}
	return nil
}

func DelSession(ctx context.Context, session string) error {
	if err := rdb.Del(ctx, redisKey(string(redisSession), session)).Err(); err != nil {
		return errors.Wrapf(err, fmt.Sprintf("failed to delete session %q", session))
	}
	return nil
}

// Used when user want to change token.
// TODO: use this whenever you are implementing server session.
func DeactivateAccessToken(ctx context.Context, token string) error {
	if err := rdb.SAdd(ctx, redisKey(string(redisExpiredAccessToken)), token).Err(); err != nil {
		return errors.Wrapf(err, fmt.Sprintf("failed to deactivate access token %q", token))
	}
	return nil
}

// 程序初始化时调用该函数初始化数据库连接
func Init() error {
	// 连接 Postgres 数据库
	var err error
	db, err = gorm.Open(postgres.Open(config.GetPostgresDsn()), &gorm.Config{})
	if err != nil || db == nil {
		return errors.Wrapf(err, "Failed to connect to Postgres")
	}
	slog.Info("Successfully initialize Postgres")

	// TODO: 连接 Redis 缓存层
	// rdb = redis.NewClient(&redis.Options{
	// 	Addr:     config.RedisAddr,
	// 	Password: config.RedisPassword,
	// 	DB:       config.RedisDB,
	// })
	// if err := rdb.Ping(context.Background()).Err(); err != nil {
	// 	slog.Error("Failed to connect Redis")
	// 	return err
	// }
	// slog.Info("Successfully initialize Redis")
	return nil
}

// 程序退出时用于
// 返回的第一个错误用于指示 Postgres 返回的错误，
// 第二个用于指示 Redis 返回的错误
func Close() (error, error) {
	pgHandle, err := db.DB()
	if err != nil {
		return errors.Wrapf(err, "failed to obtain a Postgres database handle"), nil
	}
	return func() error {
		err := pgHandle.Close()
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to close db: %s", err.Error()))
			return err
		}
		slog.Info("Postgres connection has been closed")
		return nil
	}(), nil

	// TODO: redis close error
	// rdbCloseError := func() error {
	// 	err := rdb.Close()
	// 	if err != nil {
	// 		slog.Error(fmt.Sprintf("Failed to close rdb: %s", err.Error()))
	// 		return err
	// 	}
	// 	slog.Info("Redis connection has been closed")
	// 	return nil
	// }()
}

// Migrate 自动迁移模型
func Migrate() error {
	return db.AutoMigrate(
		&User{},
		&Room{},
		&Password{},
		&Order{},
		// 用户与管理员、房间关系表
		&Relationship{},
		&WorkOrder{},
	)
}
