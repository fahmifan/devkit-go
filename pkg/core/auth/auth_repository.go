package auth

import (
	"context"

	"github.com/fahmifan/devkit/pkg/dbmodel"
	"github.com/fahmifan/devkit/pkg/xsqlc"
	"gorm.io/gorm"
)

type AuthReader struct{}

func (AuthReader) FindUserByEmail(ctx context.Context, tx *gorm.DB, email string) (authUser AuthUser, password CipherPassword, err error) {
	userModel := dbmodel.User{}
	err = tx.WithContext(ctx).Where("email = ?", email).Take(&userModel).Error
	if err != nil {
		return AuthUser{}, "", err
	}

	authUser = AuthUser{
		UserID: userModel.ID,
		Email:  userModel.Email,
		Name:   userModel.Name,
		Role:   Role(userModel.Role),
	}
	return authUser, CipherPassword(userModel.Password), nil
}

type UserWriter struct{}

func (UserWriter) SaveUser(ctx context.Context, tx xsqlc.DBTX, user *User) error {
	query := xsqlc.New(tx)

	_, err := query.SaveUser(ctx, xsqlc.SaveUserParams{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.PasswordHash,
		Role:     string(RoleUser),
	})

	return err
}
