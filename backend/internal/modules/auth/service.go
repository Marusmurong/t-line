package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	apperrors "github.com/t-line/backend/internal/pkg/errors"
	"github.com/t-line/backend/internal/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type WeChatAuth interface {
	Code2Session(code string) (openID, sessionKey string, err error)
}

type SMSSender interface {
	SendCode(phone string) error
	VerifyCode(phone, code string) (bool, error)
}

type Service struct {
	repo   *Repository
	rdb    *redis.Client
	jwtMgr *jwt.Manager
	wechat WeChatAuth
	sms    SMSSender
}

func NewService(repo *Repository, rdb *redis.Client, jwtMgr *jwt.Manager, wechat WeChatAuth, sms SMSSender) *Service {
	return &Service{
		repo:   repo,
		rdb:    rdb,
		jwtMgr: jwtMgr,
		wechat: wechat,
		sms:    sms,
	}
}

func (s *Service) WeChatLogin(ctx context.Context, req WeChatLoginReq) (*LoginResp, error) {
	openID, _, err := s.wechat.Code2Session(req.Code)
	if err != nil {
		return nil, apperrors.ErrWeChatLoginFail
	}

	user, err := s.repo.GetUserByOpenID(ctx, openID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrInternal
		}
		// create new user
		user = &User{
			WeChatOpenID: &openID,
			Nickname:     "微信用户",
			Role:         "user",
			Status:       1,
		}
		if err := s.repo.CreateUser(ctx, user); err != nil {
			return nil, fmt.Errorf("create user: %w", err)
		}
		// create wallet and points
		s.initUserAccounts(ctx, user.ID)
	}

	return s.buildLoginResp(user)
}

func (s *Service) PhoneLogin(ctx context.Context, req PhoneLoginReq) (*LoginResp, error) {
	ok, err := s.sms.VerifyCode(req.Phone, req.Code)
	if err != nil || !ok {
		return nil, apperrors.ErrSMSCodeInvalid
	}

	user, err := s.repo.GetUserByPhone(ctx, req.Phone)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrInternal
		}
		user = &User{
			Phone:    &req.Phone,
			Nickname: "用户" + req.Phone[7:],
			Role:     "user",
			Status:   1,
		}
		if err := s.repo.CreateUser(ctx, user); err != nil {
			return nil, fmt.Errorf("create user: %w", err)
		}
		s.initUserAccounts(ctx, user.ID)
	}

	return s.buildLoginResp(user)
}

func (s *Service) PasswordLogin(ctx context.Context, req PasswordLoginReq) (*LoginResp, error) {
	user, err := s.repo.GetUserByPhone(ctx, req.Phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrPasswordWrong
		}
		return nil, apperrors.ErrInternal
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apperrors.ErrPasswordWrong
	}

	return s.buildLoginResp(user)
}

func (s *Service) SendSMSCode(ctx context.Context, phone string) error {
	return s.sms.SendCode(phone)
}

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*jwt.TokenPair, error) {
	claims, err := s.jwtMgr.ParseToken(refreshToken)
	if err != nil {
		return nil, apperrors.ErrTokenExpired
	}

	// Only accept refresh tokens for token refresh
	if claims.TokenType != jwt.TokenTypeRefresh {
		return nil, apperrors.ErrUnauthorized
	}

	user, err := s.repo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, apperrors.ErrUnauthorized
	}

	return s.jwtMgr.GenerateTokenPair(user.ID, user.Role, user.MemberLevel)
}

func (s *Service) GetProfile(ctx context.Context, userID int64) (*UserResp, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, apperrors.ErrRecordNotFound
	}
	resp := ToUserResp(user)
	return &resp, nil
}

func (s *Service) UpdateProfile(ctx context.Context, userID int64, req UpdateProfileReq) (*UserResp, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, apperrors.ErrRecordNotFound
	}

	if req.Nickname != nil {
		user.Nickname = *req.Nickname
	}
	if req.AvatarURL != nil {
		user.AvatarURL = *req.AvatarURL
	}
	if req.Gender != nil {
		user.Gender = *req.Gender
	}
	if req.Age != nil {
		user.Age = *req.Age
	}
	if req.BallAge != nil {
		user.BallAge = *req.BallAge
	}
	if req.SelfLevel != nil {
		user.SelfLevel = *req.SelfLevel
	}
	if req.UTRImage != nil {
		user.UTRImage = *req.UTRImage
	}

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, apperrors.ErrInternal
	}

	resp := ToUserResp(user)
	return &resp, nil
}

func (s *Service) GetWallet(ctx context.Context, userID int64) (*WalletResp, error) {
	wallet, err := s.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.ErrRecordNotFound
	}
	return &WalletResp{
		Balance:        wallet.Balance.StringFixed(2),
		FrozenAmount:   wallet.FrozenAmount.StringFixed(2),
		TotalRecharged: wallet.TotalRecharged.StringFixed(2),
	}, nil
}

// helpers

func (s *Service) initUserAccounts(ctx context.Context, userID int64) {
	_ = s.repo.CreateWallet(ctx, &Wallet{
		UserID:  userID,
		Balance: decimal.Zero,
	})
	_ = s.repo.CreatePoints(ctx, &Points{
		UserID: userID,
	})
}

func (s *Service) buildLoginResp(user *User) (*LoginResp, error) {
	if user.Status == 0 {
		return nil, apperrors.ErrUserDisabled
	}

	tokenPair, err := s.jwtMgr.GenerateTokenPair(user.ID, user.Role, user.MemberLevel)
	if err != nil {
		return nil, apperrors.ErrInternal
	}

	return &LoginResp{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		User:         ToUserResp(user),
	}, nil
}
