package cache

import (
	context "context"
	"fmt"
	e "github.com/VSKrivoshein/test/internal/app/custom_err"
	"github.com/VSKrivoshein/test/internal/app/service"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"os"
	"time"
)

const (
	cacheKey = "users"
)

type rdb struct {
	client *redis.Client
}

func New() service.Rdb {
	adr := fmt.Sprintf("%v:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	log.Info(adr)
	rdb := &rdb{
		client: redis.NewClient(&redis.Options{
			Addr: adr,
		}),
	}
	ctx := context.Background()
	_, err := rdb.client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("redis connection ping failed: %v", err)
	}

	rdb.InvalidateCache()

	return rdb
}

func (r *rdb) GetAllUsers() ([]string, error) {
	ctx := context.Background()
	userList, err := r.client.LRange(ctx, cacheKey, 0, -1).Result()
	if err != nil {
		return nil, e.New(err, e.ErrUnexpected, codes.Internal)
	}
	return userList, nil
}

func (r *rdb) SetAllUsers(users []string) {
	ctx := context.Background()

	// подготавливаем данные, чтобы отправить их одним меньшим количеством команд
	args := make([]interface{}, 0, len(users))
	for _, user := range users {
		args = append(args, user)
	}

	var err error
	for i := 0; i < 100; i++ {
		_, err = r.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			if err := pipe.Del(ctx, cacheKey).Err(); err != nil {
				return err
			}
			if err := pipe.RPush(ctx, cacheKey, args...).Err(); err != nil {
				return err
			}
			if err := pipe.Expire(ctx, cacheKey, time.Minute).Err(); err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			continue
		}

		return
	}

	log.Fatalf("unable perform SetAllUsers err: %v", err)
}

func (r *rdb) InvalidateCache() {
	ctx := context.Background()
	var err error
	for i := 0; i < 100; i++ {
		if err = r.client.Del(ctx, cacheKey).Err(); err != nil {
			continue
		}
		return
	}

	log.Fatalf("unable invalidate cache: %v", err)
}
