package services

import (
	"chat-be/models"
	"chat-be/storage"
	"chat-be/utils"
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Signup(mobile string, username string, password string, nickname string) (string, error) {
	if mobile == "" || username == "" || password == "" || nickname == "" || !utils.IsAllNumber(mobile) || len(mobile) != 8 {
		return "", errors.New("Invalid info")
	}

	userDoc := storage.ClientDatabase.Collection("users")

	// check if exist mobile
	count, err := userDoc.CountDocuments(context.Background(), bson.D{primitive.E{Key: "mobile", Value: mobile}})
	if err != nil {
		return "", errors.New("")
	}
	if count > 0 {
		return "", errors.New("Mobile exist")
	}

	// check if exist username
	count, err = userDoc.CountDocuments(context.Background(), bson.D{primitive.E{Key: "username", Value: username}})
	if err != nil {
		return "", errors.New("")
	}
	if count > 0 {
		return "", errors.New("Username exist")
	}

	// generate token
	userToken, err := utils.GenerateToken(username + "_" + mobile)
	if err != nil {
		return "", errors.New("")
	}

	// generate password
	passwordByte := []byte(password)
	passwordEncoded, err := bcrypt.GenerateFromPassword(passwordByte, bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("")
	}

	userTemplate := models.User{
		Username:  username,
		Mobile:    mobile,
		Password:  passwordEncoded,
		Nickname:  nickname,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// add user
	result, err := userDoc.InsertOne(context.Background(), userTemplate)
	if err != nil {
		return "", errors.New("")
	}

	if result.InsertedID != nil {
		return userToken, nil
	}

	return "", errors.New("")
}

func Login(param string, password string) (string, error) {
	if param == "" || password == "" {
		return "", errors.New("Invalid info")
	}

	userDoc := storage.ClientDatabase.Collection("users")
	var userInfo models.User

	cur := userDoc.FindOne(context.Background(), bson.M{"$or": []bson.M{
		{"mobile": param},
		{"username": param},
	}})
	if cur.Err() == nil {
		return "", errors.New("Invalid Login Info")
	}
	err := cur.Decode(&userInfo)
	if err != nil {
		return "", errors.New("")
	}

	passwordByte := []byte(password)
	if bcrypt.CompareHashAndPassword(userInfo.Password, passwordByte) != nil {
		return "", errors.New("Invalid Login Info")
	}

	token, err := utils.GenerateToken(userInfo.Username + "_" + userInfo.Mobile)
	if err != nil {
		return "", errors.New("")
	}

	return token, nil
}

type getUserInfoRes struct {
	UserData models.User `json:"userData"`
	Token    string      `json:"token"`
}

func GetUserInfoByToken(token string) (getUserInfoRes, error) {
	blankData := getUserInfoRes{
		UserData: models.User{},
		Token:    "",
	}
	parsedStr, err := utils.DecodeToken(token)
	if err != nil {
		return blankData, errors.New("")
	}
	lastInd := strings.LastIndex(parsedStr, "_")
	if lastInd == -1 {
		return blankData, errors.New("")
	}

	username, mobile := parsedStr[:lastInd], parsedStr[lastInd+1:]

	userDataCur := storage.ClientDatabase.Collection("users").FindOne(context.Background(), []bson.M{{"mobile": mobile}, {"username": username}})
	if userDataCur.Err() != nil {
		return blankData, errors.New("")
	}
	var userData models.User
	err = userDataCur.Decode(&userData)
	if err != nil {
		return blankData, errors.New("")
	}

	// refresh token
	newToken, err := utils.GenerateToken(userData.Username + "_" + userData.Mobile)
	if err != nil {
		return blankData, errors.New("")
	}

	return getUserInfoRes{UserData: userData, Token: newToken}, nil
}
