package user

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/nagokos/connefut_backend/db"
	"github.com/nagokos/connefut_backend/ent"
	"github.com/nagokos/connefut_backend/ent/user"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/logger"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var host = "mailhog:1025"

type User struct {
	ID                              string
	Name                            string
	Email                           string
	Password                        string
	PasswordDigest                  string
	EmailVerificationStatus         bool
	EmailVerificationToken          string
	EmailVerificationTokenExpiresAt time.Time
}

// ** valdation **
func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(
			&u.Name,
			validation.Required.Error("名前を入力してください"),
			validation.RuneLength(1, 50).Error("名前は50文字以内で入力してください"),
		),
		validation.Field(
			&u.Email,
			validation.Required.Error("メールアドレスを入力してください"),
			is.Email.Error("メールアドレスを正しく入力してください"),
			validation.Match(regexp.MustCompile(`^[a-zA-Z0-9_+-]+(.[a-zA-Z0-9_+-]+)*@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\.)+[a-zA-Z]{2,}$`)),
		),
		validation.Field(
			&u.Password,
			validation.Required.Error("パスワードを入力してください"),
			validation.RuneLength(8, 100).Error("パスワードは8文字以上で入力してください"),
			validation.Match(regexp.MustCompile("[a-z]")).Error("パスワードを正しく入力してください"),
			validation.Match(regexp.MustCompile(`\d`)).Error("パスワードを正しく入力してください"),
		),
	)
}

func (u *User) HashGenerate() string {
	b := []byte(u.Password)
	hash, err := bcrypt.GenerateFromPassword(b, 12)
	if err != nil {
		logger.Log.Err(err)
	}
	return string(hash)
}

func (u *User) CreateUser(client *ent.UserClient, ctx context.Context) (user *model.User, err error) {
	var resUser model.User
	var pwdHash string

	if u.Password != "" {
		pwdHash = u.HashGenerate()
	}

	res, err := client.
		Create().
		SetName(u.Name).
		SetEmail(u.Email).
		SetPasswordDigest(pwdHash).
		Save(ctx)

	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintln(err))
		return &resUser, err
	}

	resUser = model.User{
		ID:    res.ID,
		Name:  res.Name,
		Email: res.Email,
	}

	return &resUser, nil
}
