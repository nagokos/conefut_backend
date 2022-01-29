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
	jwt "github.com/golang-jwt/jwt"
	"github.com/nagokos/connefut_backend/db"
	"github.com/nagokos/connefut_backend/ent"
	"github.com/nagokos/connefut_backend/ent/user"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/utils"
	"github.com/nagokos/connefut_backend/logger"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var (
	host      = "mailhog:1025"
	resUser   model.User
	SecretKey = []byte("secretKey")
)

type User struct {
	ID                              string
	Name                            string
	Email                           string
	Password                        string
	EmailVerificationStatus         bool
	EmailVerificationToken          string
	EmailVerificationTokenExpiresAt time.Time
}

// ** validation **
func (u User) CreateUserValidate() error {
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

func (u User) AuthenticateUserValidate() error {
	return validation.ValidateStruct(&u,
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

// ** utils **
func (u *User) HashGenerate() string {
	b := []byte(u.Password)
	hash, err := bcrypt.GenerateFromPassword(b, 12)
	if err != nil {
		logger.Log.Err(err)
	}
	return string(hash)
}

func (u *User) HashCompare(passwordDigest string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordDigest), []byte(u.Password))
}

func (u *User) GenerateEmailVerificationToken() string {
	h := md5.New()
	h.Write([]byte(strings.ToLower(u.Email)))
	return hex.EncodeToString(h.Sum(nil))
}

func SendVerifyEmail(emailToken string) error {
	verifyURL := fmt.Sprintf("http://localhost:8080/accounts/email_verification/%s", emailToken)
	message := strings.NewReader(verifyURL)
	transformer := japanese.ISO2022JP.NewEncoder()
	newMessage, _ := ioutil.ReadAll(transform.NewReader(message, transformer))
	err := smtp.SendMail(host, nil, "connefut@example.com", []string{"connefut@example.com"}, newMessage)
	return err
}

// ** データベース伴う処理 **
func (u *User) CreateUser(client *ent.UserClient, ctx context.Context) (user *model.User, err error) {
	pwdHash := u.HashGenerate()
	emailToken := u.GenerateEmailVerificationToken()
	tokenExpiresAt := time.Now().Add(24 * time.Hour)

	res, err := client.
		Create().
		SetName(u.Name).
		SetEmail(u.Email).
		SetPasswordDigest(pwdHash).
		SetEmailVerificationToken(emailToken).
		SetEmailVerificationTokenExpiresAt(tokenExpiresAt).
		Save(ctx)

	if err != nil {
		logger.Log.Error().Msg(err.Error())
		utils.NewValidationError("email", "このメールアドレスは既に使用されています").AddGraphQLError(ctx)
		return &resUser, err
	}

	resUser = model.User{
		ID:    res.ID,
		Name:  res.Name,
		Email: res.Email,
	}

	SendVerifyEmail(emailToken)

	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintln("fail send email: ", err))
		return &resUser, err
	}

	return &resUser, nil
}

func (u *User) AuthenticateUser(client *ent.UserClient, ctx context.Context) (*model.User, error) {
	res, err := client.
		Query().
		Where(user.Email(u.Email)).
		Only(ctx)
	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintf("user not found: %s", err))
		return &resUser, nil
	}

	err = u.HashCompare(res.PasswordDigest)
	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintf("password is incorrect: %s", err))
		utils.NewAuthenticationErorr("パスワードが正しくありません", utils.WithField("password")).AddGraphQLError(ctx)
		return &resUser, err
	}

	res, err = client.
		UpdateOneID(res.ID).
		SetLastSignInAt(time.Now()).
		Save(ctx)
	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintf("user update error: %s", err))
		return &resUser, err
	}

	resUser = model.User{
		ID:                      res.ID,
		Name:                    res.Name,
		Avatar:                  res.Avatar,
		Email:                   res.Email,
		EmailVerificationStatus: res.EmailVerificationStatus,
	}

	return &resUser, nil
}

// ** メール認証 **
func EmailVerification(w http.ResponseWriter, r *http.Request) {
	client := db.DatabaseConnection()
	defer client.Close()

	ctx := context.Background()

	token := chi.URLParam(r, "token")
	if token == "" {
		w.Write([]byte("無効なURLです"))
		logger.Log.Error().Msg("Invalid URL")
		return
	}

	res, err := client.User.
		Query().
		Where(user.EmailVerificationToken(token)).
		Only(ctx)

	if err != nil {
		w.Write([]byte("ユーザーが見つかりませんでした"))
		logger.Log.Error().Msg(fmt.Sprintf("user not found: %s", err))
		return
	}

	if time.Now().After(res.EmailVerificationTokenExpiresAt) {
		w.Write([]byte("有効期限が切れています"))
		logger.Log.Error().Msg("token expires at expired")
		return
	}

	_, err = client.User.
		UpdateOneID(res.ID).
		SetEmailVerificationStatus(true).
		SetEmailVerificationToken("").
		Save(ctx)

	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintf("user update error: %s", err))
		return
	}

	http.Redirect(w, r, "http://localhost:3000/", http.StatusMovedPermanently)
}
