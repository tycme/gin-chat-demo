package helper

import (
	"crypto/md5"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserClaims struct {
	Identity primitive.ObjectID `json:"identity"`
	Email    string             `json:"email"`
	jwt.StandardClaims
}

func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

var myKey = []byte("im")

func GenerateToken(identity, email string) (string, error) {
	objectId, err := primitive.ObjectIDFromHex(identity)
	if err != nil {
		return "", err
	}
	UserClaim := &UserClaims{
		Identity:       objectId,
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
