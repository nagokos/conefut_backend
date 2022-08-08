package user

import (
	"context"
	crand "crypto/rand"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"math/rand"
	"net/smtp"
	"regexp"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/db"
	"github.com/nagokos/connefut_backend/graph/cookie"
	"github.com/nagokos/connefut_backend/graph/jwt"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/logger"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var (
	host = "mailhog:1025"
)

var UserCtxKey = &contextKey{name: "secret"}

type contextKey struct {
	name string
}

// todo 各inputに対応したstructに
type User struct {
	Name     string
	Email    string
	Password string
}

type ChangePasswordInput struct {
	CurrentPassword         string
	NewPassword             string
	NewPasswordConfirmation string
}

type VerifyEmailInput struct {
	Code string
}

//* アドレスが重複しないかチェック
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
			err = errors.New("このメールアドレスは既に使用されています")
		}

		return err
	}
}

//* 新規パスワードと新規パスワード確認が等しいか
func passwordEqualToThePasswordConfirmation(new string) validation.RuleFunc {
	return func(value interface{}) error {
		confirmation, _ := value.(string)
		if new != confirmation {
			return errors.New("新規パスワード確認が一致しません")
		}
		return nil
	}
}

//* 新規登録
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

//* ログイン
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

//* 新規メールアドレス
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

//* 認証コード
func (i VerifyEmailInput) VerifyEmailValidate() error {
	return validation.ValidateStruct(&i,
		validation.Field(
			&i.Code,
			validation.Required.Error("認証コードを入力してください"),
			validation.Match(regexp.MustCompile(`^[0-9]{6}$`)).Error("認証コードに誤りがあります"),
		),
	)
}

//* パスワード変更
func (i ChangePasswordInput) ChangePasswordValidate() error {
	return validation.ValidateStruct(&i,
		validation.Field(
			&i.CurrentPassword,
			validation.Required.Error("現在のパスワードを入力してください"),
			validation.RuneLength(8, 100).Error("現在のパスワードは8文字以上で入力してください"),
			validation.Match(regexp.MustCompile("[a-z]")).Error("現在のパスワードを正しく入力してください"),
			validation.Match(regexp.MustCompile(`\d`)).Error("現在のパスワードを正しく入力してください"),
		),
		validation.Field(
			&i.NewPassword,
			validation.Required.Error("新規パスワードを入力してください"),
			validation.RuneLength(8, 100).Error("新規パスワードは8文字以上で入力してください"),
			validation.Match(regexp.MustCompile("[a-z]")).Error("新規パスワードを正しく入力してください"),
			validation.Match(regexp.MustCompile(`\d`)).Error("新規パスワードを正しく入力してください"),
		),
		validation.Field(
			&i.NewPasswordConfirmation,
			validation.Required.Error("新規パスワード確認を入力してください"),
			validation.By(passwordEqualToThePasswordConfirmation(i.NewPassword)),
		),
	)
}

//* ログインユーザー取得
func GetViewer(ctx context.Context) *model.User {
	raw, _ := ctx.Value(UserCtxKey).(*model.User)
	return raw
}

//* パスワードのハッシュを生成
func GeneratePasswordHash(password string) string {
	b := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(b, 12)
	if err != nil {
		logger.NewLogger().Error(err.Error())
	}
	return string(hash)
}

//* ユーザーのハッシュ化したパスワードと送られてきたパスワードを比較
func CheckPasswordHash(passwordDigest, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordDigest), []byte(password))
}

//* メール認証のPINを生成
func GenerateEmailVerification() (string, error) {
	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	rand.Seed(seed.Int64())
	var pin string
	for i := 0; i < 6; i++ {
		pin = fmt.Sprintf(pin+"%v", rand.Intn(9))
	}
	return pin, nil
}

//* 実際にメールを送信する処理
func SendingVerifyEmail(pin string, to string) error {
	message := strings.NewReader(fmt.Sprint(pin))
	transformer := japanese.ISO2022JP.NewEncoder()
	newMessage, _ := ioutil.ReadAll(transform.NewReader(message, transformer))
	err := smtp.SendMail(host, nil, "connefut@example.com", []string{to}, newMessage)
	return err
}

//* idからユーザーを取得
func GetUser(ctx context.Context, dbPool *pgxpool.Pool, id int) (*model.User, error) {
	cmd := "SELECT id, name, email, avatar, introduction, email_verification_status, unverified_email FROM users WHERE id = $1"

	var user model.User
	row := dbPool.QueryRow(ctx, cmd, id)
	err := row.Scan(&user.DatabaseID, &user.Name, &user.Email, &user.Avatar,
		&user.Introduction, &user.EmailVerificationStatus, &user.UnverifiedEmail)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return &user, nil
}

//* メールアドレスに認証メール送信
func SendVerifyEmail(ctx context.Context, dbPool *pgxpool.Pool) (bool, error) {
	viewer := GetViewer(ctx)
	pinExpiresAt := time.Now().Add(10 * time.Minute)
	pin, err := GenerateEmailVerification()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return false, err
	}

	cmd := `
	  UPDATE users
		SET (email_verification_pin, email_verification_pin_expires_at) = ($1, $2)
		WHERE id = $3
	`
	if _, err := dbPool.Exec(
		ctx, cmd,
		pin, pinExpiresAt, viewer.DatabaseID,
	); err != nil {
		logger.NewLogger().Error(err.Error())
		return false, err
	}

	if err := SendingVerifyEmail(pin, *viewer.UnverifiedEmail); err != nil {
		logger.NewLogger().Error(err.Error())
		return false, err
	}
	return true, err
}

//* 新しいメールアドレスに認証メール送信
func (u User) SendVerifyNewEmail(ctx context.Context, dbPool *pgxpool.Pool) (model.SendVerifyNewEmailResult, error) {
	cmd := `
	  UPDATE users
		SET (email_verification_pin, email_verification_pin_expires_at, unverified_email) = ($1, $2, $3)
		WHERE id = $4
		RETURNING id, unverified_email
	`
	viewer := GetViewer(ctx)
	pinExpiresAt := time.Now().Add(10 * time.Minute)
	pin, err := GenerateEmailVerification()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	row := dbPool.QueryRow(
		ctx, cmd,
		pin, pinExpiresAt, u.Email, viewer.DatabaseID,
	)
	var user model.User
	if err := row.Scan(&user.DatabaseID, &user.UnverifiedEmail); err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	if err := SendingVerifyEmail(pin, *user.UnverifiedEmail); err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	result := model.SendVerifyNewEmailSuccess{Viewer: &model.Viewer{AccountUser: &user}}
	return result, nil
}

//* メールアドレス認証
func (i VerifyEmailInput) VerifyEmail(ctx context.Context, dbPool *pgxpool.Pool) (model.VerifyEmailResult, error) {
	viewer := GetViewer(ctx)
	cmd := `
	  SELECT email_verification_pin, email_verification_pin_expires_at
		FROM users
		WHERE id = $1
	`
	row := dbPool.QueryRow(ctx, cmd, viewer.DatabaseID)
	var pin string
	var pinExpiresAt time.Time
	if err := row.Scan(&pin, &pinExpiresAt); err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	//* pinの有効期限チェック
	if time.Now().Local().After(pinExpiresAt) {
		logger.NewLogger().Error("pin code expired")
		result := model.VerifyEmailPinExpiredError{
			Message: "認証コードの有効期限が切れています。認証コードを再取得してください。",
		}
		return result, nil
	}

	//* 送られてきたpinとユーザーのpinを比較
	if i.Code != pin {
		logger.NewLogger().Error("Pin code does not match")
		result := model.VerifyEmailAuthenticationError{
			Message: "認証コードに誤りがあります",
		}
		return result, nil
	}

	cmd = `
	  UPDATE users
		SET (email, unverified_email, email_verification_status, email_verification_pin_expires_at, email_verification_pin, updated_at) = ($1, $2, $3, $4, $5, $6)
		WHERE id = $7
		RETURNING id, email, unverified_email, email_verification_status
	`
	row = dbPool.QueryRow(
		ctx, cmd,
		viewer.UnverifiedEmail, nil, "verified", nil, nil, time.Now().Local(), viewer.DatabaseID,
	)
	var user model.User
	if err := row.Scan(&user.DatabaseID, &user.Email, &user.UnverifiedEmail, &user.EmailVerificationStatus); err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	result := model.VerifyEmailSuccess{Viewer: &model.Viewer{AccountUser: &user}}
	return result, nil
}

//* パスワードを変更
func (i ChangePasswordInput) ChangePassword(ctx context.Context, dbPool *pgxpool.Pool) (model.ChangePasswordResult, error) {
	viewer := GetViewer(ctx)
	cmd := `
	  SELECT password_digest
		FROM users
		WHERE id = $1
	`
	row := dbPool.QueryRow(ctx, cmd, viewer.DatabaseID)
	var passwordDigest *string
	if err := row.Scan(&passwordDigest); err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	//* 送られてきた現在のパスワードとハッシュ化したパスワードを比較
	if err := CheckPasswordHash(*passwordDigest, i.CurrentPassword); err != nil {
		logger.NewLogger().Error(err.Error())
		result := model.ChangePasswordAuthenticationError{
			Message: "現在のパスワードが有効ではありません",
		}
		return result, nil
	}

	//* ハッシュを生成
	hash := GeneratePasswordHash(i.NewPassword)
	cmd = `
	  UPDATE users
		SET (password_digest, created_at) = ($1, $2)
		WHERE id = $3
	`
	if _, err := dbPool.Exec(
		ctx, cmd,
		hash, time.Now().Local(), viewer.DatabaseID,
	); err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	return model.ChangePasswordSuccess{IsChangedPassword: true}, nil
}

//* 外部認証用のユーザー取得
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

//* ユーザー新規登録
func (u *User) RegisterUser(ctx context.Context, dbPool *pgxpool.Pool) (model.RegisterUserResult, error) {
	cmd := `
		INSERT INTO users
			(name, email, unverified_email, password_digest, email_verification_pin,
				email_verification_pin_expires_at, last_sign_in_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, name, email, avatar, email_verification_status, introduction, unverified_email
	`
	now := time.Now().Local()
	pwdHash := GeneratePasswordHash(u.Password)
	pinExpiresAt := time.Now().Add(10 * time.Minute)
	pin, err := GenerateEmailVerification()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	row := dbPool.QueryRow(
		ctx, cmd,
		u.Name, u.Email, u.Email, pwdHash, pin, pinExpiresAt, now, now, now,
	)

	var user model.User
	err = row.Scan(&user.DatabaseID, &user.Name, &user.Email, &user.Avatar, &user.EmailVerificationStatus, &user.Introduction, &user.UnverifiedEmail)
	if err != nil {
		return nil, err
	}

	//todo 本番環境と開発環境では処理が違ってくる
	err = SendingVerifyEmail(pin, user.Email)
	if err != nil {
		return nil, err
	}

	result := model.RegisterUserSuccess{Viewer: &model.Viewer{AccountUser: &user}}
	token, err := jwt.GenerateToken(user.DatabaseID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	cookie.SetAuthCookie(ctx, token)
	return result, nil
}

func (u *User) LoginUser(ctx context.Context, dbPool *pgxpool.Pool) (model.LoginUserResult, error) {
	cmd := `
	  SELECT id, name, email, avatar, email_verification_status, introduction, unverified_email, password_digest 
		FROM users 
		WHERE email = $1
	`
	var user model.User
	var passwordDigest *string
	row := dbPool.QueryRow(ctx, cmd, u.Email)
	if err := row.Scan(&user.DatabaseID, &user.Name, &user.Email, &user.Avatar, &user.EmailVerificationStatus,
		&user.Introduction, &user.UnverifiedEmail, &passwordDigest,
	); err != nil {
		// todo レコードが見つからない場合はschemaで返すかgraphqlErrorで返すか とりあえずschemaで返す
		if err == pgx.ErrNoRows {
			logger.NewLogger().Error("user not found")
			result := model.LoginUserNotFoundError{
				Message: "メールアドレス、またはパスワードが正しくありません",
			}
			return result, nil
		} else {
			fmt.Println(err)
			logger.NewLogger().Error(err.Error())
			return nil, err
		}
	}

	if err := CheckPasswordHash(*passwordDigest, u.Password); err != nil {
		logger.NewLogger().Error(err.Error())
		result := model.LoginUserAuthenticationError{
			Message: "メールアドレス、またはパスワードが正しくありません",
		}
		return result, nil
	}

	result := model.LoginUserSuccess{Viewer: &model.Viewer{AccountUser: &user}}
	token, err := jwt.GenerateToken(user.DatabaseID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	cookie.SetAuthCookie(ctx, token)
	return result, nil
}
