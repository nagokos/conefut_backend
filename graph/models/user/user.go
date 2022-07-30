package user

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	jwt "github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/db"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/utils"
	"github.com/nagokos/connefut_backend/logger"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var (
	host      = "mailhog:1025"
	SecretKey = []byte("secretKey")
)

var UserCtxKey = &contextKey{name: "secret"}

type contextKey struct {
	name string
}

type User struct {
	ID                              string
	Name                            string
	Email                           string
	Password                        string
	EmailVerificationStatus         bool
	EmailVerificationToken          string
	EmailVerificationTokenExpiresAt time.Time
}

type NullableUser struct {
	ID     *string
	Name   *string
	Avatar *string
}

func checkExistsEmail() validation.RuleFunc {
	return func(v interface{}) error {
		var err error

		email := v.(string)
		dbPool := db.DatabaseConnection()

		cmd := "SELECT COUNT(DISTINCT id) FROM users WHERE email = $1"
		row := dbPool.QueryRow(context.Background(), cmd, email)

		var count int
		err = row.Scan(&count)

		if err != nil {
			logger.NewLogger().Error(err.Error())
			return err
		}

		if count == 1 {
			logger.NewLogger().Error("This email address is already exists")
			err = errors.New("このメールアドレスは既に存在します")
		}

		return err
	}
}

// ** validation **
func (u User) CreateUserValidate() error {
	return validation.ValidateStruct(&u,
		validation.Field(
			&u.Name,
			validation.Required.Error("名前を入力してください"),
			validation.RuneLength(1, 20).Error("名前は50文字以内で入力してください"),
		),
		validation.Field(
			&u.Email,
			validation.Required.Error("メールアドレスを入力してください"),
			validation.RuneLength(1, 100).Error("メールアドレスは100文字以内で入力してください"),
			validation.Match(regexp.MustCompile(`^[a-zA-Z0-9_+-]+(.[a-zA-Z0-9_+-]+)*@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\.)+[a-zA-Z]{2,}$`)).
				Error("メールアドレスを正しく入力してください"),
			validation.By(checkExistsEmail()),
		),
		validation.Field(
			&u.Password,
			validation.Required.Error("パスワードを入力してください"),
			validation.RuneLength(8, 100).Error("パスワードは8~100文字で入力してください"),
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
			validation.Match(regexp.MustCompile(`^[a-zA-Z0-9_+-]+(.[a-zA-Z0-9_+-]+)*@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\.)+[a-zA-Z]{2,}$`)).
				Error("メールアドレスを正しく入力してください"),
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

func (u User) SendVerifyNewEmailValidate() error {
	return validation.ValidateStruct(&u,
		validation.Field(
			&u.Email,
			validation.Required.Error("メールアドレスを入力してください"),
			validation.RuneLength(1, 100).Error("メールアドレスは100文字以内で入力してください"),
			validation.Match(regexp.MustCompile(`^[a-zA-Z0-9_+-]+(.[a-zA-Z0-9_+-]+)*@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\.)+[a-zA-Z]{2,}$`)).
				Error("メールアドレスを正しく入力してください"),
			validation.By(checkExistsEmail()),
		),
	)
}

// ** utils **
func GenerateHash(password string) string {
	b := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(b, 12)
	if err != nil {
		logger.NewLogger().Error(err.Error())
	}
	return string(hash)
}

func CheckPasswordHash(passwordDigest, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passwordDigest), []byte(password))
	return err
}

func GenerateEmailVerificationToken() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func CreateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("secretKey"))
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return "", err
	}

	return tokenString, nil
}

func SendVerifyEmail(emailToken string) error {
	verifyURL := fmt.Sprintf("http://localhost:8080/accounts/verify_email?token=%s", emailToken)
	message := strings.NewReader(verifyURL)
	transformer := japanese.ISO2022JP.NewEncoder()
	newMessage, _ := ioutil.ReadAll(transform.NewReader(message, transformer))
	err := smtp.SendMail(host, nil, "connefut@example.com", []string{"connefut@example.com"}, newMessage)
	return err
}

func GetUser(ctx context.Context, dbPool *pgxpool.Pool, id string) (*model.User, error) {
	cmd := "SELECT id, name, email, avatar, introduction, email_verification_status FROM users WHERE id = $1"

	var user model.User
	row := dbPool.QueryRow(ctx, cmd, utils.DecodeUniqueID(id))
	err := row.Scan(&user.DatabaseID, &user.Name, &user.Email, &user.Avatar, &user.Introduction, &user.EmailVerificationStatus)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return &user, nil
}

func GetUserIDByProviderAndUID(ctx context.Context, dbPool *pgxpool.Pool, provider, uid string) (int, error) {
	cmd := `
	  SELECT u.id, u.name, u.avatar, u.introduction
		FROM users as u
		INNER JOIN authentications as o
		  ON u.id = o.user_id
		WHERE u.id = (
			SELECT user_id
			FROM authentications
			WHERE provider = $1
		  AND uid = $2
		)
	`
	row := dbPool.QueryRow(
		ctx, cmd,
		provider, uid,
	)

	var user model.User
	if err := row.Scan(&user.DatabaseID, &user.Name, &user.Avatar, &user.Introduction); err != nil {
		logger.NewLogger().Error(err.Error())
		return 0, err
	}
	return user.DatabaseID, nil
}

// ** データベース伴う処理 **
func (u *User) RegisterUser(ctx context.Context, dbPool *pgxpool.Pool) (*model.RegisterUserPayload, error) {
	pwdHash := HashGenerate(u.Password)
	emailToken := u.GenerateEmailVerificationToken()
	tokenExpiresAt := time.Now().Add(24 * time.Hour)

	cmd := `
		INSERT INTO users
			(name, email, password_digest, email_verification_token, 
				email_verification_token_expires_at, last_sign_in_at, created_at, updated_at
			) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id, name, email, avatar, email_verification_status
	`

	row := dbPool.QueryRow(
		ctx, cmd,
		u.Name, u.Email, pwdHash, emailToken, tokenExpiresAt, time.Now().Local(), time.Now().Local(), time.Now().Local(),
	)

	var payload model.RegisterUserPayload
	var viewer model.User

	err := row.Scan(&viewer.DatabaseID, &viewer.Name, &viewer.Email, &viewer.Avatar, &viewer.EmailVerificationStatus)
	if err != nil {
		return nil, err
	}

	// 本番環境と開発環境では違う
	err = SendVerifyEmail(emailToken)
	if err != nil {
		return nil, err
	}

	payload.Viewer = &viewer
	return &payload, nil
}

func (u *User) LoginUser(ctx context.Context, dbPool *pgxpool.Pool) (*model.LoginUserPayload, error) {
	var payload model.LoginUserPayload
	var viewer model.User
	var passwordDigest string

	cmd := "SELECT id, name, email, avatar, email_verification_status, password_digest FROM users WHERE email = $1"
	row := dbPool.QueryRow(ctx, cmd, u.Email)

	err := row.Scan(&viewer.DatabaseID, &viewer.Name, &viewer.Email, &viewer.Avatar, &viewer.EmailVerificationStatus, &passwordDigest)
	if err != nil {
		payload.UserErrors = append(payload.UserErrors, model.LoginUserAuthenticationError{
			Message: "メールアドレス、またはパスワードが正しくありません",
		})
		return &payload, err
	}

	err = CheckPasswordHash(passwordDigest, u.Password)
	if err != nil {
		payload.UserErrors = append(payload.UserErrors, model.LoginUserAuthenticationError{
			Message: "メールアドレス、またはパスワードが正しくありません",
		})
		return &payload, err
	}

	payload.Viewer = &viewer
	return &payload, nil
}

// ** メール認証 **
func EmailVerification(w http.ResponseWriter, r *http.Request) {
	dbPool := db.DatabaseConnection()
	defer dbPool.Close()

	token := chi.URLParam(r, "token")
	if token == "" {
		_, err := w.Write([]byte("無効なURLです"))
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		logger.NewLogger().Error("Invalid URL")
		return
	}

	ctx := context.Background()

	cmd := "SELECT id, email_verification_token_expires_at FROM users WHERE email_verification_token = $1"
	row := dbPool.QueryRow(ctx, cmd, token)

	var ID int
	var tokenExpiresAt time.Time
	err := row.Scan(&ID, &tokenExpiresAt)

	if err != nil {
		_, err = w.Write([]byte("ユーザーが見つかりませんでした"))
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		logger.NewLogger().Sugar().Errorf("user not found: %s", err)
		return
	}

	if time.Now().After(tokenExpiresAt) {
		_, err = w.Write([]byte("有効期限が切れています"))
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		logger.NewLogger().Error("token expires at expired")
		return
	}

	cmd = `
	  UPDATE users AS u
		SET email_verification_status = $1, email_verification_token = $2, updated_at = $3
		WHERE u.id = $4
	`
	_, err = dbPool.Exec(ctx, cmd, "verified", nil, time.Now().Local(), ID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		_, err = w.Write([]byte("メールアドレスの認証に失敗しました"))
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		return
	}

	jwt, _ := CreateToken(ID)
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    jwt,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(time.Hour * 1),
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "http://localhost:3000/", http.StatusMovedPermanently)
}
