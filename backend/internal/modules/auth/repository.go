package auth

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// User operations

func (r *Repository) CreateUser(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *Repository) GetUserByID(ctx context.Context, id int64) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByOpenID(ctx context.Context, openID string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("wechat_openid = ?", openID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByPhone(ctx context.Context, phone string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Wallet operations

func (r *Repository) CreateWallet(ctx context.Context, wallet *Wallet) error {
	return r.db.WithContext(ctx).Create(wallet).Error
}

func (r *Repository) GetWalletByUserID(ctx context.Context, userID int64) (*Wallet, error) {
	var wallet Wallet
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *Repository) UpdateWalletWithVersion(ctx context.Context, wallet *Wallet) error {
	result := r.db.WithContext(ctx).
		Model(wallet).
		Where("id = ? AND version = ?", wallet.ID, wallet.Version).
		Updates(map[string]interface{}{
			"balance":        wallet.Balance,
			"frozen_amount":  wallet.FrozenAmount,
			"total_recharged": wallet.TotalRecharged,
			"version":        wallet.Version + 1,
		})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *Repository) CreateWalletTransaction(ctx context.Context, tx *WalletTransaction) error {
	return r.db.WithContext(ctx).Create(tx).Error
}

func (r *Repository) ListWalletTransactions(ctx context.Context, walletID int64, offset, limit int) ([]WalletTransaction, int64, error) {
	var txs []WalletTransaction
	var total int64

	query := r.db.WithContext(ctx).Where("wallet_id = ?", walletID)
	query.Model(&WalletTransaction{}).Count(&total)
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&txs).Error

	return txs, total, err
}

// Points operations

func (r *Repository) CreatePoints(ctx context.Context, points *Points) error {
	return r.db.WithContext(ctx).Create(points).Error
}

func (r *Repository) GetPointsByUserID(ctx context.Context, userID int64) (*Points, error) {
	var points Points
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&points).Error; err != nil {
		return nil, err
	}
	return &points, nil
}

// Transaction helper

func (r *Repository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
