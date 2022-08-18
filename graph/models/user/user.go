package user

import (
	"context"
	crand "crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/big"
	"math/rand"
	"net/http"
	"net/smtp"
	"regexp"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/99designs/gqlgen/graphql"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/db"
	"github.com/nagokos/connefut_backend/graph/cookie"
	"github.com/nagokos/connefut_backend/graph/jwt"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/prefecture"
	"github.com/nagokos/connefut_backend/graph/models/sport"
	"github.com/nagokos/connefut_backend/logger"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const (
	host             = "mailhog:1025"
	avatarBucket     = "connefut-user-upload"
	avatarObjectPath = "avatar/"
)

var UserCtxKey = &contextKey{name: "secret"}

type contextKey struct {
	name string
}

// todo 各inputに対応したstructに
type User struct {
	Name          string
	Email         string
	Password      string
	Introduction  string
	Avatar        graphql.Upload
	PrefectureIDs []int
	SportIDs      []int
	WebsiteURL    string
}

type ChangePasswordInput struct {
	CurrentPassword         string
	NewPassword             string
	NewPasswordConfirmation string
}

type ChangeEmailInput struct {
	NewEmail string
}

type VerifyEmailInput struct {
	Code string
}

type ResetPasswordInput struct {
	Email                   string
	NewPassword             string
	NewPasswordConfirmation string
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
			return errors.New("新規パスワードと一致しません")
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
func (i ChangeEmailInput) ChangeEmailValidate() error {
	return validation.ValidateStruct(&i,
		validation.Field(
			&i.NewEmail,
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

//* パスワードリセット
func (i ResetPasswordInput) SendResetPasswordEmailValidate() error {
	return validation.ValidateStruct(&i,
		validation.Field(
			&i.Email,
			validation.Required.Error("メールアドレスを入力してください"),
			validation.RuneLength(1, 100).Error("メールアドレスは100文字以内で入力してください"),
			validation.Match(regexp.MustCompile(`^[a-zA-Z0-9_+-]+(.[a-zA-Z0-9_+-]+)*@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\.)+[a-zA-Z]{2,}$`)).
				Error("メールアドレスを正しく入力してください"),
		),
	)
}

//* パスワードリセット変更
func (i ResetPasswordInput) ResetPasswordValidate() error {
	return validation.ValidateStruct(&i,
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

func (u User) UpdateUserValidate() error {
	return validation.ValidateStruct(&u,
		validation.Field(
			&u.Name,
			validation.Required.Error("名前を入力してください"),
			validation.RuneLength(1, 20).Error("名前は50文字以内で入力してください"),
		),
		validation.Field(
			&u.Introduction,
			validation.RuneLength(0, 160).Error("自己紹介は160文字以内で入力してください"),
		),
		validation.Field(
			&u.WebsiteURL,
			is.URL.Error("URLの形式が正しくありません"),
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

//* メール認証のCodeを生成
func GenerateEmailVerificationCode() (string, error) {
	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	rand.Seed(seed.Int64())
	var code string
	for i := 0; i < 6; i++ {
		code = fmt.Sprintf(code+"%v", rand.Intn(9))
	}
	return code, nil
}

//* パスワードリセットのトークンを生成
func GeneratePasswordResetToken() (string, error) {
	b := make([]byte, 32)
	if _, err := crand.Read(b); err != nil {
		logger.NewLogger().Error(err.Error())
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

//* 実際にメールを送信する処理
func SendingEmail(content string, to string) error {
	message := strings.NewReader(content)
	transformer := japanese.ISO2022JP.NewEncoder()
	newMessage, _ := ioutil.ReadAll(transform.NewReader(message, transformer))
	err := smtp.SendMail(host, nil, "connefut@example.com", []string{to}, newMessage)
	return err
}

//* idからユーザーを取得
func GetUser(ctx context.Context, dbPool *pgxpool.Pool, id int) (*model.User, error) {
	cmd := `
	  SELECT id, name, email, avatar, introduction, email_verification_status, unverified_email, website_url 
		FROM users 
		WHERE id = $1
	`

	var user model.User
	row := dbPool.QueryRow(ctx, cmd, id)
	err := row.Scan(&user.DatabaseID, &user.Name, &user.Email, &user.Avatar,
		&user.Introduction, &user.EmailVerificationStatus, &user.UnverifiedEmail, &user.WebsiteURL)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return &user, nil
}

//* メールアドレスに認証メール送信
func SendVerifyEmail(ctx context.Context, dbPool *pgxpool.Pool) (bool, error) {
	viewer := GetViewer(ctx)
	now := time.Now().Local()
	expiresAt := now.Add(10 * time.Minute)
	code, err := GenerateEmailVerificationCode()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return false, err
	}

	cmd := `
	  UPDATE users
		SET (email_verification_code, email_verification_code_expires_at, updated_at) = ($1, $2)
		WHERE id = $3
	`
	if _, err := dbPool.Exec(
		ctx, cmd,
		code, expiresAt, viewer.DatabaseID,
	); err != nil {
		logger.NewLogger().Error(err.Error())
		return false, err
	}

	if err := SendingEmail(code, *viewer.UnverifiedEmail); err != nil {
		logger.NewLogger().Error(err.Error())
		return false, err
	}
	return true, err
}

//* 新しいメールアドレスに認証メール送信
func (i *ChangeEmailInput) ChangeEmail(ctx context.Context, dbPool *pgxpool.Pool) (model.ChangeUserEmailResult, error) {
	cmd := `
	  UPDATE users
		SET (email_verification_code, email_verification_code_expires_at, unverified_email) = ($1, $2, $3)
		WHERE id = $4
		RETURNING id, unverified_email
	`
	viewer := GetViewer(ctx)
	codeExpiresAt := time.Now().Add(10 * time.Minute)
	code, err := GenerateEmailVerificationCode()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	row := dbPool.QueryRow(
		ctx, cmd,
		code, codeExpiresAt, i.NewEmail, viewer.DatabaseID,
	)
	var user model.User
	if err := row.Scan(&user.DatabaseID, &user.UnverifiedEmail); err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	if err := SendingEmail(code, *user.UnverifiedEmail); err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	result := model.ChangeUserEmailSuccess{Viewer: &model.Viewer{AccountUser: &user}}
	return result, nil
}

//* メールアドレス認証
func (i VerifyEmailInput) VerifyEmail(ctx context.Context, dbPool *pgxpool.Pool) (model.VerifyUserEmailResult, error) {
	viewer := GetViewer(ctx)
	cmd := `
	  SELECT email_verification_code, email_verification_code_expires_at
		FROM users
		WHERE id = $1
	`
	row := dbPool.QueryRow(ctx, cmd, viewer.DatabaseID)
	var code string
	var codeExpiresAt time.Time
	if err := row.Scan(&code, &codeExpiresAt); err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	//* codeの有効期限チェック
	if time.Now().Local().After(codeExpiresAt) {
		logger.NewLogger().Error("code code expired")
		result := model.VerifyUserEmailCodeExpiredError{
			Message: "認証コードの有効期限が切れています。認証コードを再取得してください。",
		}
		return result, nil
	}

	//* 送られてきたcodeとユーザーのcodeを比較
	if i.Code != code {
		logger.NewLogger().Error("Code does not match")
		result := model.VerifyUserEmailAuthenticationError{
			Message: "認証コードに誤りがあります",
		}
		return result, nil
	}

	cmd = `
	  UPDATE users
		SET (email, unverified_email, email_verification_status, email_verification_code_expires_at, email_verification_code, updated_at) = ($1, $2, $3, $4, $5, $6)
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
	result := model.VerifyUserEmailSuccess{Viewer: &model.Viewer{AccountUser: &user}}
	return result, nil
}

//* パスワードを変更
func (i ChangePasswordInput) ChangePassword(ctx context.Context, dbPool *pgxpool.Pool) (model.ChangeUserPasswordResult, error) {
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
		result := model.ChangeUserPasswordAuthenticationError{
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
	return model.ChangeUserPasswordSuccess{IsChangedPassword: true}, nil
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
			(name, email, unverified_email, password_digest, email_verification_code,
				email_verification_code_expires_at, last_sign_in_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, name, email, avatar, email_verification_status, introduction, unverified_email
	`
	now := time.Now().Local()
	pwdHash := GeneratePasswordHash(u.Password)
	codeExpiresAt := time.Now().Add(10 * time.Minute)
	code, err := GenerateEmailVerificationCode()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	row := dbPool.QueryRow(
		ctx, cmd,
		u.Name, u.Email, u.Email, pwdHash, code, codeExpiresAt, now, now, now,
	)

	var user model.User
	err = row.Scan(&user.DatabaseID, &user.Name, &user.Email, &user.Avatar, &user.EmailVerificationStatus, &user.Introduction, &user.UnverifiedEmail)
	if err != nil {
		return nil, err
	}

	//todo 本番環境と開発環境では処理が違ってくる
	err = SendingEmail(code, user.Email)
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

//* ログイン
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

//* リセットパスワードの送信
func (i *ResetPasswordInput) SendResetPasswordEmail(ctx context.Context, dbPool *pgxpool.Pool) (model.SendResetPasswordEmailToUserResult, error) {
	cmd := `
	  SELECT COUNT(DISTINCT id)
		FROM users
		WHERE email = $1
	`
	row := dbPool.QueryRow(ctx, cmd, i.Email)
	var count int
	if err := row.Scan(&count); err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	if count > 0 {
		cmd := `
		  UPDATE users
			SET (password_reset_token, password_reset_token_expires_at, updated_at) = ($1, $2, $3)
			WHERE email = $4
		`
		now := time.Now().Local()
		expiresAt := now.Add(1 * time.Hour)
		token, err := GeneratePasswordResetToken()
		if err != nil {
			logger.NewLogger().Error(err.Error())
			return nil, err
		}
		if _, err := dbPool.Exec(ctx, cmd, token, expiresAt, now, i.Email); err != nil {
			logger.NewLogger().Error(err.Error())
			return nil, err
		}
		resetPasswordURL := fmt.Sprintf("http://localhost:8080/password/reset?token=%s", token)
		if err := SendingEmail(resetPasswordURL, i.Email); err != nil {
			logger.NewLogger().Error(err.Error())
			return nil, err
		}
		result := model.SendResetPasswordEmailToUserSuccess{
			IsSentEmail: true,
		}
		return result, nil
	} else {
		result := model.SendResetPasswordEmailToUserNotFoundError{
			Message: "ユーザーが見つかりませんでした",
		}
		return result, nil
	}
}

//* パスワードリセットトークンの有効性を確認
func IsTokenValid(ctx context.Context, dbPool *pgxpool.Pool, token string) (bool, error) {
	cmd := `
	  SELECT password_reset_token_expires_at
		FROM users
		WHERE password_reset_token = $1
	`
	row := dbPool.QueryRow(ctx, cmd, token)
	var expiresAt time.Time
	if err := row.Scan(&expiresAt); err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	if time.Now().Local().After(expiresAt) {
		logger.NewLogger().Error("Token has expired")
		return false, nil
	}
	return true, nil
}

//* パスワードリセットURLの有効性の確認
func ConfirmationPasswordResetURL(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		logger.NewLogger().Error("there is no token")
		http.Error(w, "there is no token", http.StatusBadRequest)
		return
	}

	dbPool := db.DatabaseConnection()
	defer dbPool.Close()
	ctx := context.Background()

	isValid, err := IsTokenValid(ctx, dbPool, token)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		http.Error(w, "token is invalid", http.StatusBadRequest)
		return
	}

	if isValid {
		redirect := fmt.Sprintf("http://localhost:5173/account/password_reset?token=%s", token)
		http.Redirect(w, r, redirect, http.StatusPermanentRedirect)
	} else {
		http.Redirect(w, r, "http://localhost:5173/login", http.StatusPermanentRedirect)
	}
}

func (i *ResetPasswordInput) ResetPassword(ctx context.Context, dbPool *pgxpool.Pool, token string) (model.ResetUserPasswordResult, error) {
	isValid, err := IsTokenValid(ctx, dbPool, token)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	if isValid {
		cmd := `
		  UPDATE users
			SET (password_digest, password_reset_token, password_reset_token_expires_at, updated_at) = ($1, $2, $3, $4)
			WHERE password_reset_token = $5
			RETURNING id, name, email, avatar, email_verification_status, introduction, unverified_email
		`
		pwdhash := GeneratePasswordHash(i.NewPassword)
		now := time.Now().Local()
		row := dbPool.QueryRow(ctx, cmd, pwdhash, nil, nil, now, token)
		var user model.User
		if err := row.Scan(&user.DatabaseID, &user.Name, &user.Email, &user.Avatar,
			&user.EmailVerificationStatus, &user.Introduction, &user.UnverifiedEmail,
		); err != nil {
			logger.NewLogger().Error(err.Error())
			return nil, err
		}
		result := model.ResetUserPasswordSuccess{
			Viewer: &model.Viewer{
				AccountUser: &user,
			},
		}
		jwt, err := jwt.GenerateToken(user.DatabaseID)
		if err != nil {
			logger.NewLogger().Error(err.Error())
			return nil, err
		}
		cookie.SetAuthCookie(ctx, jwt)
		return result, nil
	} else {
		logger.NewLogger().Error("invalid token")
		result := model.ResetUserPasswordInvalidTokenError{
			Message: "トークンが無効です",
		}
		return result, nil
	}
}

//* ユーザー情報更新
func (u *User) UpdateUser(ctx context.Context, dbPool *pgxpool.Pool) (model.UpdateUserResult, error) {
	cmd := `
	  UPDATE users
		SET (name, introduction, website_url, updated_at) = ($1, $2, $3, $4)
		WHERE id = $5
		RETURNING id, name, introduction, website_url
	`
	now := time.Now().Local()
	viewer := GetViewer(ctx)
	row := dbPool.QueryRow(ctx, cmd, u.Name, u.Introduction, u.WebsiteURL, now, viewer.DatabaseID)
	var user model.User
	if err := row.Scan(&user.DatabaseID, &user.Name, &user.Introduction, &user.WebsiteURL); err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	//* 現在の活動エリアを取得
	currentPrefectures, err := prefecture.GetPrefecturesByUserID(ctx, dbPool, viewer.DatabaseID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	//* 削除する都道府県をを取得
	var removePrefectures []int
	for _, currentPrefecture := range currentPrefectures {
		found := false
		for _, sentPrefecture := range u.PrefectureIDs {
			if currentPrefecture.DatabaseID == sentPrefecture {
				found = true
			}
		}
		if !found {
			removePrefectures = append(removePrefectures, currentPrefecture.DatabaseID)
		}
	}
	//* 付与する都道府県を取得
	var addPrefectures []int
	for _, sentPrefecture := range u.PrefectureIDs {
		found := false
		for _, currentPrefecture := range currentPrefectures {
			if sentPrefecture == currentPrefecture.DatabaseID {
				found = true
			}
		}
		if !found {
			addPrefectures = append(addPrefectures, sentPrefecture)
		}
	}

	//* 活動エリアのトランザクション開始
	tx, err := dbPool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	//* 活動エリアを付与
	for _, addPrefecture := range addPrefectures {
		if err := prefecture.AddUserActivityArea(ctx, tx, viewer.DatabaseID, addPrefecture); err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}
	}
	//* 活動エリアを削除
	for _, removePrefecture := range removePrefectures {
		if err := prefecture.RemoveUserActivieArea(ctx, tx, viewer.DatabaseID, removePrefecture); err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}
	}
	//* 活動エリアの付与、削除が成功したらコミットする
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	//* 現在のプレイスポーツを取得
	currentSports, err := sport.GetSportsByUserID(ctx, dbPool, viewer.DatabaseID)
	if err != nil {
		return nil, err
	}
	//* 削除するスポーツを取得
	var removeSports []int
	for _, currentSport := range currentSports {
		found := false
		for _, sentSport := range u.SportIDs {
			if currentSport.DatabaseID == sentSport {
				found = true
			}
		}
		if !found {
			removeSports = append(removeSports, currentSport.DatabaseID)
		}
	}
	//* 付与するスポーツを取得
	var addSports []int
	for _, sentSport := range u.SportIDs {
		found := false
		for _, currentSport := range currentSports {
			if sentSport == currentSport.DatabaseID {
				found = true
			}
		}
		if !found {
			addSports = append(addSports, sentSport)
		}
	}

	//* スポーツのトランザクションを開始
	tx, err = dbPool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	//* スポーツを付与
	for _, addSport := range addSports {
		if err := sport.AddUserPlaySport(ctx, tx, viewer.DatabaseID, addSport); err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}
	}
	//* スポーツを削除
	for _, removeSport := range removeSports {
		if err := sport.RemoveUserPlaySport(ctx, tx, viewer.DatabaseID, removeSport); err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}
	}
	//* スポーツの付与、削除に成功したらコミット
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	result := model.UpdateUserSuccess{
		Viewer: &model.Viewer{
			AccountUser: &user,
		},
	}
	return result, nil
}

//* ユーザーアバターアップロード
func (u *User) UploadUserAvatar(ctx context.Context, dbPool *pgxpool.Pool, gcsClient *storage.Client) (model.UploadUserAvatarResult, error) {
	viewer := GetViewer(ctx)
	// todo 別パッケージにアップロード系は切り出す
	encodeObject := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf("%d:%s", viewer.DatabaseID, u.Avatar.Filename)))
	extension := strings.Split(u.Avatar.Filename, ".")[5]
	object := fmt.Sprintf("%s.%s", encodeObject, extension)

	writer := gcsClient.Bucket(avatarBucket).Object(avatarObjectPath + object).NewWriter(ctx)
	if _, err := io.Copy(writer, u.Avatar.File); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	cmd := `
	UPDATE users
	SET (avatar, updated_at) = ($1, $2)
	WHERE id = $3
	RETURNING id, avatar
	`
	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s%s/%s", avatarBucket, writer.Bucket, writer.Name)
	now := time.Now().Local()
	row := dbPool.QueryRow(ctx, cmd, publicURL, now, viewer.DatabaseID)
	var user model.User
	if err := row.Scan(&user.DatabaseID, &user.Avatar); err != nil {
		return nil, err
	}

	if publicURL != viewer.Avatar {
		deleteObject := viewer.Avatar[strings.Index(viewer.Avatar, avatarObjectPath):]
		o := gcsClient.Bucket(avatarBucket).Object(deleteObject)
		exists, _ := o.Attrs(ctx)
		if exists != nil {
			if err := o.Delete(ctx); err != nil {
				return nil, err
			}
		}
	}

	return &model.UploadUserAvatarSuccess{Viewer: &model.Viewer{AccountUser: &user}}, nil
}
