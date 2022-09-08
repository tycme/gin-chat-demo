package helper

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/dgrijalva/jwt-go"
	"github.com/jordan-wright/email"
	"github.com/tycme/gin-chat/define"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

var myKey = []byte("im")

func GenerateToken(identity, email string) (string, error) {
	UserClaim := &UserClaims{
		Identity:       identity,
		Email:          email,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func AnalyseToken(tokenString string) (*UserClaims, error) {
	userClaims := &UserClaims{}
	claims, err := jwt.ParseWithClaims(tokenString, userClaims, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("analyse Token error: %v", err)
	}
	return userClaims, nil
}

func SendCode(toUserEmail, code string) error {
	e := email.NewEmail()
	e.From = "Get <honort@163.com>"
	e.To = []string{toUserEmail}
	e.Subject = "验证码已发送, 请查收"
	e.HTML = []byte("您的验证码： <b>" + code + "<b>")
	return e.SendWithTLS("smtp.163.com:465",
		smtp.PlainAuth("", "honort@163.com", define.MailPassword, "smtp.163.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"},
	)
}
