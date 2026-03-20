package sms

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/t-line/backend/internal/pkg/logger"
)

const (
	codeExpiry    = 5 * time.Minute
	codeCooldown  = 60 * time.Second
	codeKeyPrefix = "sms:code:"
	coolKeyPrefix = "sms:cool:"
)

type Sender struct {
	rdb      *redis.Client
	provider string
	// TODO: add actual SMS provider client (aliyun/tencent)
}

func NewSender(rdb *redis.Client, provider string) *Sender {
	return &Sender{
		rdb:      rdb,
		provider: provider,
	}
}

func (s *Sender) SendCode(phone string) error {
	ctx := context.Background()

	// check cooldown
	coolKey := coolKeyPrefix + phone
	if s.rdb.Exists(ctx, coolKey).Val() > 0 {
		return fmt.Errorf("发送过于频繁")
	}

	code := generateCode()

	// store code in Redis
	codeKey := codeKeyPrefix + phone
	if err := s.rdb.Set(ctx, codeKey, code, codeExpiry).Err(); err != nil {
		return fmt.Errorf("store code: %w", err)
	}

	// set cooldown
	s.rdb.Set(ctx, coolKey, "1", codeCooldown)

	// TODO: call actual SMS provider to send code
	// Only log masked phone in debug mode; never log the code in production
	if os.Getenv("APP_MODE") == "debug" {
		logger.L.Debugf("[SMS] code sent to %s***%s", phone[:3], phone[len(phone)-4:])
	}

	return nil
}

func (s *Sender) VerifyCode(phone, code string) (bool, error) {
	ctx := context.Background()
	codeKey := codeKeyPrefix + phone

	stored, err := s.rdb.Get(ctx, codeKey).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("get code: %w", err)
	}

	if stored != code {
		return false, nil
	}

	// delete code after verification
	s.rdb.Del(ctx, codeKey)
	return true, nil
}

func generateCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
