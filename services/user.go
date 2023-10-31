package services

import (
	"chat-be/models"
	"chat-be/storage"
	"chat-be/utils"
	"context"
	"errors"
	"fmt"
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

	// generate password
	passwordByte := []byte(password)
	passwordEncoded, err := bcrypt.GenerateFromPassword(passwordByte, bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("")
	}

	userTemplate := models.User{
		Username:   username,
		Mobile:     mobile,
		Password:   passwordEncoded,
		Nickname:   nickname,
		Status:     "active",
		LastActive: time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// add user
	result, err := userDoc.InsertOne(context.Background(), userTemplate)
	if err != nil {
		return "", errors.New("")
	}

	if result.InsertedID != nil {

		// generate token
		userToken, err := utils.GenerateToken(result.InsertedID.(primitive.ObjectID).Hex())
		if err != nil {
			return "", errors.New("")
		}
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

	token, err := utils.GenerateToken(userInfo.Id.Hex())
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
	newToken, err := utils.GenerateToken(userData.Id.Hex())
	if err != nil {
		return blankData, errors.New("")
	}

	err = UpdateUserInfo(userData.Id.String(), models.User{})
	if err != nil {
		return blankData, errors.New("")
	}

	return getUserInfoRes{UserData: userData, Token: newToken}, nil
}

func UpdateUserInfo(_id string, userinfo models.User) error {
	if _id == "" {
		return errors.New("Invalid Info")
	}
	userId := utils.ToObjectId(_id)

	userDoc := storage.ClientDatabase.Collection("users")

	updatedField := models.User{}
	if userinfo.Nickname != "" {
		updatedField.Nickname = userinfo.Nickname
		updatedField.UpdatedAt = time.Now()
	}
	if userinfo.Avatar != "" {
		updatedField.Avatar = userinfo.Avatar
		updatedField.UpdatedAt = time.Now()
	}
	if userinfo.Describe != "" {
		updatedField.Describe = userinfo.Describe
		updatedField.UpdatedAt = time.Now()
	}
	updatedField.LastActive = time.Now()

	result, err := userDoc.UpdateByID(context.Background(), userId, updatedField)
	if err != nil {
		return errors.New("")
	}
	fmt.Println("Updated: ", result.ModifiedCount, updatedField)
	return nil
}

func ChangePassword(_id string, oldPass string, newPass string) error {
	if _id == "" || oldPass == "" || newPass == "" {
		return errors.New("Invalid Info")
	}
	userDoc := storage.ClientDatabase.Collection("users")
	userCur := userDoc.FindOne(context.Background(), bson.M{"_id": utils.ToObjectId(_id)})
	if userCur.Err() != nil {
		return errors.New("Invalid Info")
	}

	var userInfo models.User
	err := userCur.Decode(&userInfo)
	if err != nil {
		return errors.New("")
	}

	oldPassByte := []byte(oldPass)
	if bcrypt.CompareHashAndPassword(userInfo.Password, oldPassByte) != nil {
		return errors.New("Invalid Password")
	}

	// generate password
	newPassByte := []byte(newPass)
	passwordEncoded, err := bcrypt.GenerateFromPassword(newPassByte, bcrypt.DefaultCost)
	if err != nil {
		return errors.New("")
	}

	_, err = userDoc.UpdateByID(context.Background(), utils.ToObjectId(_id), bson.M{"password": passwordEncoded})
	if err != nil {
		return errors.New("")
	}

	return nil
}
