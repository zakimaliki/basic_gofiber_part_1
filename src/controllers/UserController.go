package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"gofiber/src/helpers"
	"gofiber/src/models"
	"gofiber/src/services"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// func RegisterUser(c *fiber.Ctx) error {
// 	var user models.User
// 	if err := c.BodyParser(&user); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Failed to parse request body",
// 		})
// 	}

// 	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
// 	user.Password = string(hashPassword)

// 	models.PostUser(&user)

// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
// 		"message": "Register successfully",
// 	})
// }

func RegisterUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	// var existingUser models.User
	// if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
	// 	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
	// 		"error": "Email already used",
	// 	})
	// }

	existingUser := models.FindEmail(&user)
	if len(existingUser) > 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Email already existing",
		})
	}

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashPassword)
	user.Verify = "false"

	randomID, err := models.GenerateRandom6DigitID()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate user ID",
		})
	}
	user.ID = randomID

	uuid := strconv.Itoa(int(randomID))

	tokenBytes := make([]byte, 64)
	rand.Read(tokenBytes)
	token := hex.EncodeToString(tokenBytes)

	url := os.Getenv("BASE_URL") + "verify?id=" + uuid + "&token=" + token

	if err := services.SendEmail(user.Email, "Verify Email", url); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send verification email",
		})
	}

	if err := models.PostUser(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	userVerification := models.UserVerification{
		UserID: user.ID,
		Token:  token,
	}

	if err := models.PostUserVerify(&userVerification); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user verification",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Sign Up Success, Please check your email for verification",
	})
}

func VerifyAccount(c *fiber.Ctx) error {
	queryUsersId := c.Query("id")
	queryToken := c.Query("token")

	if queryUsersId == "" || queryToken == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Invalid url verification",
		})
	}

	userID, err := strconv.ParseUint(queryUsersId, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	user := models.FindID(uint(userID))
	if user.ID == 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Error users has not found",
		})
	}

	if user.Verify != "false" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Users has been verified",
		})
	}

	userVer, err := models.CheckUsersVerification(uint(userID), queryToken)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Error invalid credential verification",
		})
	}

	if err := models.UpdateUserVerify(int(userID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update account verification",
		})
	}

	if err := models.DeleteUsersVerification(userVer.ID, queryToken); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user verification",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Users verified successfully",
	})
}

func LoginUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	validateEmail := models.FindEmail(&user)
	if len(validateEmail) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Email not found",
		})
	}

	if validateEmail[0].Verify == "false" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User is unverify",
		})
	}

	var PasswordSecond string
	for _, user := range validateEmail {
		PasswordSecond = user.Password
	}

	if err := bcrypt.CompareHashAndPassword([]byte(PasswordSecond), []byte(user.Password)); err != nil || user.Password == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid password",
		})
	}

	jwtKey := os.Getenv("SECRETKEY")
	payload := map[string]interface{}{
		"email": user.Email,
	}

	token, err := helpers.GenerateToken(jwtKey, payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not generate access token",
		})
	}

	refreshToken, err := helpers.GenerateRefreshToken(jwtKey, payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not generate refresh token",
		})
	}

	item := map[string]string{
		"Email":        user.Email,
		"Token":        token,
		"RefreshToken": refreshToken,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Login successfully",
		"data":    item,
	})
}

func RefreshToken(c *fiber.Ctx) error {
	var input struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	jwtKey := os.Getenv("SECRETKEY")
	payload := map[string]interface{}{
		"refreshToken": input.RefreshToken,
	}

	token, err := helpers.GenerateToken(jwtKey, payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not generate access token",
		})
	}

	refreshToken, err := helpers.GenerateRefreshToken(jwtKey, payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not generate refresh token",
		})
	}
	item := map[string]string{
		"Token":        token,
		"RefreshToken": refreshToken,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Login successfully",
		"data":    item,
	})
}
