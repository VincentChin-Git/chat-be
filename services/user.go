package services

import (
	"chat-be/models"
	"chat-be/storage"
	"chat-be/utils"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Signup(mobile string, username string, password string) (string, error) {
	if mobile == "" || username == "" || password == "" || !utils.IsAllNumber(mobile) || len(mobile) != 8 {
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

	nickname, err := utils.GenerateRandomNumber(5)
	if err != nil {
		return "", errors.New("")
	}
	nickname = "user-" + nickname

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
	userId, err := utils.DecodeToken(token)
	if err != nil {
		return blankData, errors.New("")
	}

	userDataCur := storage.ClientDatabase.Collection("users").FindOne(context.Background(), []bson.M{{"_id": utils.ToObjectId(userId)}})
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

	_, err = userDoc.UpdateByID(context.Background(), utils.ToObjectId(_id), bson.M{
		"password":  passwordEncoded,
		"updatedAt": time.Now(),
	})
	if err != nil {
		return errors.New("")
	}

	return nil
}

func SearchUser(key string, value string) (models.User, error) {
	userDoc := storage.ClientDatabase.Collection("user")
	cur := userDoc.FindOne(context.Background(), bson.M{key: value, "status": "active"})
	if cur.Err() != nil {
		fmt.Println(cur.Err())
		return models.User{}, nil
	}

	var userTemp, userRes models.User

	err := cur.Decode(&userTemp)
	if err != nil {
		fmt.Println(err.Error())
		return models.User{}, errors.New("")
	}

	userRes.Avatar = userTemp.Avatar
	userRes.Describe = userTemp.Describe
	userRes.Nickname = userTemp.Nickname
	userRes.LastActive = userTemp.LastActive
	userRes.Id = userTemp.Id
	userRes.Mobile = userTemp.Mobile

	return userRes, nil
}

func ForgetPassword(_id string, code string, password string) error {
	if _id == "" || code == "" || password == "" {
		return errors.New("Invalid Info")
	}

	// search for matched code
	resetPassDoc := storage.ClientDatabase.Collection("resetPass")
	result := resetPassDoc.FindOne(context.Background(), bson.M{
		"_id":    utils.ToObjectId(_id),
		"code":   code,
		"status": "pending",
	})

	if result.Err() != nil {
		fmt.Println(result.Err().Error())
		return errors.New("")
	}

	// generate password
	newPassByte := []byte(password)
	passwordEncoded, err := bcrypt.GenerateFromPassword(newPassByte, bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return errors.New("")
	}

	var resetInfo models.ResetPass
	err = result.Decode(&resetInfo)
	if err != nil {
		fmt.Println(err)
		return errors.New("")
	}

	// get userid and update user password
	userId := resetInfo.UserId

	userDoc := storage.ClientDatabase.Collection("users")
	_, err = userDoc.UpdateByID(context.Background(), utils.ToObjectId(userId.Hex()), bson.M{
		"password":  passwordEncoded,
		"updatedAt": time.Now(),
	})

	if err != nil {
		fmt.Println(err)
		return errors.New("")
	}

	// update reset password to completed
	_, err = resetPassDoc.UpdateByID(context.Background(), utils.ToObjectId(_id), bson.M{
		"status":    "completed",
		"updatedAt": time.Now(),
	})

	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func AddForgetPassword(userId string, code string) error {
	if userId == "" || code == "" {
		return errors.New("Invalid Info")
	}

	resetPassDoc := storage.ClientDatabase.Collection("resetPass")

	getCode := storage.ReadRedis(userId + "_forgetPassword")
	if getCode == code {
		resetPassInfo := models.ResetPass{
			UserId:     utils.ToObjectId(userId),
			VerifyCode: code,
			Status:     "pending",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		_, err := resetPassDoc.InsertOne(context.Background(), resetPassInfo)
		if err != nil {
			return errors.New("")
		}

		return nil
	} else {
		return errors.New("Invalid Verification Code")
	}
}

func GetForgetPassCode(mobile string) (string, error) {
	if mobile == "" {
		return "", errors.New("Invalid Info")
	}

	userDoc := storage.ClientDatabase.Collection("users")
	cur := userDoc.FindOne(context.Background(), bson.M{"mobile": mobile})
	if cur.Err() != nil {
		fmt.Println(cur.Err().Error())
		return "", errors.New("No User Found")
	}
	var userInfo models.User
	err := cur.Decode(&userInfo)
	if err != nil {
		fmt.Println(err.Error())
		return "", errors.New("")
	}

	generatedCode, err := utils.GenerateRandomNumber(6)
	if err != nil {
		fmt.Println(err.Error())
		return "", errors.New("")
	}

	isErr := storage.WriteRedis(userInfo.Id.Hex()+"_forgetPassword", generatedCode, time.Minute)

	if isErr {
		return "", errors.New("")
	}

	return userInfo.Id.Hex(), nil
}
