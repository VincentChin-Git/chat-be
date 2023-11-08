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
	"go.mongodb.org/mongo-driver/mongo"
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

	timeNow := time.Now()
	userTemplate := models.User{
		Username:   username,
		Mobile:     mobile,
		Password:   passwordEncoded,
		Nickname:   nickname,
		Status:     "active",
		LastActive: &timeNow,
		CreatedAt:  &timeNow,
		UpdatedAt:  &timeNow,
	}

	// add user
	result, err := userDoc.InsertOne(context.Background(), userTemplate)
	if err != nil {
		return "", errors.New("")
	}

	if result.InsertedID != nil {

		// generate token
		userToken, err := utils.GenerateToken(result.InsertedID.(primitive.ObjectID).Hex(), time.Now().Add(30*24*time.Hour))
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
	if cur.Err() != nil {
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

	token, err := utils.GenerateToken(userInfo.Id.Hex(), time.Now().Add(30*24*time.Hour))
	if err != nil {
		return "", errors.New("")
	}

	return token, nil
}

type getUserInfoRes struct {
	UserData models.User `json:"userData"`
	Token    string      `json:"token"`
}

func GetUserInfoByToken(userId string) (getUserInfoRes, error) {
	blankData := getUserInfoRes{
		UserData: models.User{},
		Token:    "",
	}

	fmt.Println(userId, "userId")

	userDataCur := storage.ClientDatabase.Collection("users").FindOne(context.Background(), bson.M{"_id": utils.ToObjectId(userId)})
	if userDataCur.Err() != nil {
		return blankData, errors.New("")
	}
	var userData models.User
	err := userDataCur.Decode(&userData)
	if err != nil {
		return blankData, errors.New("")
	}

	// refresh token
	newToken, err := utils.GenerateToken(userData.Id.Hex(), time.Now().Add(30*24*time.Hour))
	if err != nil {
		return blankData, errors.New("")
	}

	err = UpdateUserInfo(userData.Id.Hex(), models.User{})
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
	timeNow := time.Now()
	if userinfo.Nickname != "" {
		updatedField.Nickname = userinfo.Nickname
		updatedField.UpdatedAt = &timeNow
	}
	if userinfo.Avatar != "" {
		updatedField.Avatar = userinfo.Avatar
		updatedField.UpdatedAt = &timeNow
	}
	if userinfo.Describe != "" {
		updatedField.Describe = userinfo.Describe
		updatedField.UpdatedAt = &timeNow
	}
	updatedField.LastActive = &timeNow

	result, err := userDoc.UpdateByID(context.Background(), userId, bson.M{"$set": updatedField})
	if err != nil {
		fmt.Println(err.Error())
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
		"$set": bson.M{
			"password":  passwordEncoded,
			"updatedAt": time.Now(),
		},
	})
	if err != nil {
		return errors.New("")
	}

	return nil
}

func SearchUser(key string, value string, userId string) (models.User, error) {
	userDoc := storage.ClientDatabase.Collection("users")
	var cur *mongo.SingleResult
	if key == "_id" {
		cur = userDoc.FindOne(context.Background(), bson.M{key: utils.ToObjectId(value), "status": "active"})
	} else {
		cur = userDoc.FindOne(context.Background(), bson.M{key: value, "status": "active"})
	}
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
	if userTemp.Id.Hex() == userId {
		return models.User{}, nil
	}

	userRes.Avatar = userTemp.Avatar
	userRes.Describe = userTemp.Describe
	userRes.Nickname = userTemp.Nickname
	userRes.LastActive = userTemp.LastActive
	userRes.Id = userTemp.Id
	userRes.Mobile = userTemp.Mobile

	return userRes, nil
}

func ForgetPassword(userId string, code string, password string) error {
	if userId == "" || code == "" || password == "" {
		return errors.New("Invalid Info")
	}

	fmt.Println(userId, code, password)

	// search for matched code
	resetPassDoc := storage.ClientDatabase.Collection("resetPass")
	result := resetPassDoc.FindOne(context.Background(), bson.M{
		"userId":     utils.ToObjectId(userId),
		"verifyCode": code,
		"status":     "pending",
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

	userDoc := storage.ClientDatabase.Collection("users")
	_, err = userDoc.UpdateByID(context.Background(), utils.ToObjectId(userId), bson.M{
		"$set": bson.M{
			"password":  passwordEncoded,
			"updatedAt": time.Now(),
		},
	})

	if err != nil {
		fmt.Println(err)
		return errors.New("")
	}

	// update reset password to completed
	_, err = resetPassDoc.UpdateByID(context.Background(), utils.ToObjectId(resetInfo.Id.Hex()), bson.M{
		"$set": bson.M{
			"status":    "completed",
			"updatedAt": time.Now(),
		},
	})

	if err != nil {
		fmt.Println(err)
		return errors.New("")
	}

	return nil
}

func AddForgetPassword(userId string, code string) error {
	if userId == "" || code == "" {
		return errors.New("Invalid Info")
	}

	resetPassDoc := storage.ClientDatabase.Collection("resetPass")

	getCode := storage.ReadRedis(userId + "_forgetPassword")
	timeNow := time.Now()
	userIdObject := utils.ToObjectId(userId)
	if getCode == code {
		resetPassInfo := models.ResetPass{
			UserId:     &userIdObject,
			VerifyCode: code,
			Status:     "pending",
			CreatedAt:  &timeNow,
			UpdatedAt:  &timeNow,
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
