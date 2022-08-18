package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/stheven26/db"
	"github.com/stheven26/globals"
	"github.com/stheven26/helpers"
	"github.com/stheven26/models"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
			"time":    time.Now(),
		})
	}

	hash, err := helpers.HashPassword(data["password"])

	if err != nil {
		return c.JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
			"time":    time.Now(),
		})
	}

	user := models.User{
		Username: data["username"],
		Password: hash,
	}

	connection := db.GetConnectionDB()
	connection.Create(&user)
	return c.JSON(fiber.Map{
		"status": http.StatusOK,
		"data":   user,
		"time":   time.Now(),
	})
}

func Login(c *fiber.Ctx) error {
	var data map[string]string
	var user models.User

	if err := c.BodyParser(&data); err != nil {
		return c.JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
			"time":    time.Now(),
		})
	}

	connection := db.GetConnectionDB()
	connection.Where(`username = ?`, data["username"]).First(&user)

	if user.Id == 0 {
		return c.JSON(fiber.Map{
			"status":  http.StatusNotFound,
			"message": "User not found",
			"time":    time.Now(),
		})
	}

	_, err := helpers.CheckHashPassword(data["password"], user.Password)

	if err != nil {
		return c.JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "password incorrect",
			"time":    time.Now(),
		})
	}

	//generate JWT token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(72 * time.Hour).Unix(),
	})

	token, err := claims.SignedString([]byte(globals.Key))

	if err != nil {
		return c.JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
			"time":    time.Now(),
		})
	}

	//Generate HMAC 256 from JWT Token
	h := hmac.New(sha256.New, []byte(globals.Key))

	h.Write([]byte(token))

	hmac := hex.EncodeToString(h.Sum(nil))

	cookies := fiber.Cookie{
		Name:     "jwt",
		Value:    hmac,
		Expires:  time.Now().Add(72 * time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookies)

	return c.JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Success Login",
		"time":    time.Now(),
	})
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Success Logout",
		"time":    time.Now(),
	})
}

func GetAllBlog(c *fiber.Ctx) error {
	var data []models.Blog

	connection := db.GetConnectionDB()
	connection.Find(&data)

	return c.JSON(fiber.Map{
		"status": http.StatusOK,
		"data":   data,
		"time":   time.Now(),
	})
}

func GetBlogById(c *fiber.Ctx) error {
	var data models.Blog

	connection := db.GetConnectionDB()

	if err := connection.Where(`id = ?`, c.Params("id")).First(&data).Error; err != nil {
		return c.JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
			"time":    time.Now(),
		})
	}

	return c.JSON(fiber.Map{
		"status": http.StatusOK,
		"data":   data,
		"time":   time.Now(),
	})
}

func CreateBlog(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(globals.Key), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthorized",
			"time":    time.Now(),
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	connection := db.GetConnectionDB()
	connection.Where("id", claims.Issuer).First(&user)

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
			"time":    time.Now(),
		})
	}

	blog := models.Blog{
		Title: data["title"],
		Body:  data["body"],
		Slug:  data["slug"],
	}

	connection.Create(&blog)

	return c.JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": blog,
		"time":    time.Now(),
	})
}

func UpdateBlog(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(globals.Key), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthorized",
			"time":    time.Now(),
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	connection := db.GetConnectionDB()
	connection.Where("id", claims.Issuer).First(&user)

	var data map[string]string
	var blog models.Blog

	if err := connection.Where(`id`, c.Params("id")).First(&blog).Error; err != nil {
		return c.JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Blog tidak ditemukan",
			"time":    time.Now(),
		})
	}

	if err := c.BodyParser(&data); err != nil {
		return c.JSON(fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
			"time":    time.Now(),
		})
	}

	updateData := models.Blog{
		Title: data["title"],
		Body:  data["body"],
		Slug:  data["slug"],
	}

	connection.Model(&blog).Updates(&updateData)

	return c.JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Berhasil Update data",
		"data":    data,
		"time":    time.Now(),
	})
}

func DeleteBlog(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(globals.Key), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthorized",
			"time":    time.Now(),
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	connection := db.GetConnectionDB()
	connection.Where("id", claims.Issuer).First(&user)

	var blog models.Blog

	if err := connection.Where(`id = ?`, c.Params("id")).First(&blog).Error; err != nil {
		return c.JSON(fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "data blog tidak ditemukan",
			"time":    time.Now(),
		})
	}

	connection.Delete(&blog)

	return c.JSON(fiber.Map{
		"status":  http.StatusOK,
		"message": "Berhasil menghapus blog",
		"time":    time.Now(),
	})
}
