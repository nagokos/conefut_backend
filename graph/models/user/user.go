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
			err = errors.New("このメールアドレスは既に存在します")
		}

		return err
	}
}

//* 新規パスワードと新規パスワード確認が等しいか
func passwordEqualToThePasswordConfirmation(new string) validation.RuleFunc {
	return func(value interface{}) error {
		confirmation, _ := value.(string)
		if new != confirmation {
			return errors.New("新規パスワードと新規パスワード確認が一致しません")
		}
		return nil
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

func (i VerifyEmailInput) VerifyEmailValidate() error {
	return validation.ValidateStruct(&i,
		validation.Field(
			&i.Code,
			validation.Required.Error("認証コードを入力してください"),
			validation.Match(regexp.MustCompile(`^[0-9]{6}$`)).Error("認証コードに誤りがあります"),
		),
	)
}

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
func GenerateHash(password string) string {
	b := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(b, 12)
	if err != nil {
		logger.NewLogger().Error(err.Error())
	}
	return string(hash)
}

//* ユーザーのハッシュ化したパスワードと送られてきたパスワードを比較
func CheckPasswordHash(passwordDigest, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passwordDigest), []byte(password))
	return err
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

//* Cookieにセットする認証トークンを生成(JWT)
func CreateToken(userID int) (string, error) {
	now := time.Now().Local()
	payload := jwt.MapClaims{
		"sub": userID,
		"exp": now.Add(time.Hour * 24).Unix(),
		"iat": now.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	tokenString, err := token.SignedString([]byte("secretKey"))
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return "", err
	}

	return tokenString, nil
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
func GetUser(ctx context.Context, dbPool *pgxpool.Pool, id string) (*model.User, error) {
	cmd := "SELECT id, name, email, avatar, introduction, email_verification_status, unverified_email FROM users WHERE id = $1"

	var user model.User
	row := dbPool.QueryRow(ctx, cmd, utils.DecodeUniqueID(id))
	err := row.Scan(&user.DatabaseID, &user.Name, &user.Email, &user.Avatar,
		&user.Introduction, &user.EmailVerificationStatus, &user.UnverifiedEmail)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return &user, nil
}

func GetViewer(ctx context.Context) *model.User {
	raw, _ := ctx.Value(UserCtxKey).(*model.User)
	return raw
}

func ReSendVerifyEmail(ctx context.Context, dbPool *pgxpool.Pool) (bool, error) {
	viewer := GetViewer(ctx)
	tokenExpiresAt := time.Now().Local().Add(24 * time.Hour)
	emailToken, err := GenerateEmailVerificationToken()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return false, err
	}

	cmd := `
	  UPDATE users
		SET (email_verification_token, email_verification_token_expires_at) = ($1, $2)
		WHERE id = $3
	`
	if _, err := dbPool.Exec(
		ctx, cmd,
		emailToken, tokenExpiresAt, viewer.DatabaseID,
	); err != nil {
		logger.NewLogger().Error(err.Error())
		return false, nil
	}

	if err := SendVerifyEmail(emailToken); err != nil {
		logger.NewLogger().Error(err.Error())
		return false, nil
	}
	return true, err
}

func SendVerifyNewEmail(ctx context.Context, dbPool *pgxpool.Pool, email string) (*model.SendVerifyNewEmailPayload, error) {
	var payload model.SendVerifyNewEmailPayload

	viewer := GetViewer(ctx)
	tokenExpiresAt := time.Now().Local().Add(24 * time.Hour)
	emailToken, err := GenerateEmailVerificationToken()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	cmd := `
	  UPDATE users
		SET (email_verification_token, email_verification_token_expires_at) = ($1, $2)
		WHERE id = $3
	`
	if _, err := dbPool.Exec(
		ctx, cmd,
		emailToken, tokenExpiresAt, viewer.DatabaseID,
	); err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	verifyURL := fmt.Sprintf("http://localhost:8080/accounts/verify_new_email?token=%s&email=%s", emailToken, base64.URLEncoding.EncodeToString([]byte(email)))
	message := strings.NewReader(verifyURL)
	transformer := japanese.ISO2022JP.NewEncoder()
	newMessage, _ := ioutil.ReadAll(transform.NewReader(message, transformer))
	if err := smtp.SendMail(host, nil, "connefut@example.com", []string{"connefut@example.com"}, newMessage); err != nil {
		return nil, err
	}
	payload.IsSendVerifyEmail = true

	return &payload, nil
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
	pwdHash := GenerateHash(u.Password)
	tokenExpiresAt := time.Now().Add(24 * time.Hour)
	emailToken, err := GenerateEmailVerificationToken()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

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

	err = row.Scan(&viewer.DatabaseID, &viewer.Name, &viewer.Email, &viewer.Avatar, &viewer.EmailVerificationStatus)
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
func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	dbPool := db.DatabaseConnection()
	defer dbPool.Close()

	token := r.URL.Query().Get("token")
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
		SET (email_verification_status, email_verification_token, email_verification_token_expires_at, updated_at) = ($1, $2, $3, $4)
		WHERE u.id = $5
	`
	_, err = dbPool.Exec(ctx, cmd, "verified", nil, nil, time.Now().Local(), ID)
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
	http.Redirect(w, r, os.Getenv("CLIENT_BASE_URL"), http.StatusMovedPermanently)
}

func VerifyNewEmail(w http.ResponseWriter, r *http.Request) {
	dbPool := db.DatabaseConnection()

	encodedEmail := r.URL.Query().Get("email")
	if encodedEmail == "" {
		logger.NewLogger().Error("email not found")
		http.Error(w, "email not found", http.StatusBadRequest)
		return
	}
	email, err := base64.URLEncoding.DecodeString(encodedEmail)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cmd := `
	  SELECT COUNT(DISTINCT id)
		FROM users
		WHERE email = $1
	`
	row := dbPool.QueryRow(
		r.Context(), cmd,
		string(email),
	)
	var count int
	if err := row.Scan(&count); err != nil {
		logger.NewLogger().Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count > 0 {
		logger.NewLogger().Error("This email address is already exists")
		http.Error(w, "This email address is already exists", http.StatusBadRequest)
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		logger.NewLogger().Error("token not found")
		http.Error(w, "token not found", http.StatusBadRequest)
		return
	}

	cmd = `
	  SELECT id, email_verification_token_expires_at
		FROM users
		WHERE email_verification_token = $1
	`
	row = dbPool.QueryRow(
		r.Context(), cmd,
		token,
	)

	var userID int
	var tokenExpiresAt time.Time
	if err := row.Scan(&userID, &tokenExpiresAt); err != nil {
		logger.NewLogger().Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if time.Now().After(tokenExpiresAt) {
		logger.NewLogger().Error("token expires at expired")
		http.Error(w, "token expires at expired", http.StatusBadRequest)
		return
	}

	cmd = `
	  UPDATE users
		SET (email, email_verification_token, email_verification_token_expires_at, updated_at) = ($1, $2, $3, $4)
		WHERE id = $5
	`
	if _, err := dbPool.Exec(
		r.Context(), cmd,
		string(email), nil, nil, time.Now().Local(), userID,
	); err != nil {
		logger.NewLogger().Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jwt, err := CreateToken(userID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    jwt,
		HttpOnly: true,
		Path:     "/",
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(time.Hour * 24),
	})
	http.Redirect(w, r, os.Getenv("CLIENT_BASE_URL"), http.StatusMovedPermanently)
}
