package controller

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var secretKey = []byte("blog_project")

type User struct {
	Id       uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Email    string    `json:"email"`
}

type Register_User struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	CheckPassword string `json:"checkpassword"`
	Email         string `json:"email"`
}

type Login_User struct {
	Email    string `json:"Email"`
	Password string `json:"password"`
}

type CustomClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type ChangePassword struct {
	UserID          string `json:"userId"`
	OldPassword     string `json:"oldPassword"`
	NewPassword     string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}

func GetUser(c *gin.Context) {
	db := connectDB()
	var users []*User
	db.Find(&users)
	closeDB(db)
	c.JSON(200, users)
}

func GetUserById(c *gin.Context) {
	db := connectDB()
	var user User
	queryResult := db.Where("id = ?", c.Param("id")).First(&user)
	if queryResult.Error != nil {
		c.JSON(500, gin.H{
			"message": "query error" + queryResult.Error.Error(),
		})
		closeDB(db)
		return
	}
	closeDB(db)
	c.JSON(200, user)
}

func isValidEmail(email string) bool {
	// 定義電子信箱格式的正則表達式
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// 使用正則表達式檢查電子信箱格式
	match, err := regexp.MatchString(emailRegex, email)
	if err != nil {
		fmt.Println("Error checking email format:", err)
		return false
	}

	return match
}

func HandleRegister(c *gin.Context) {
	var rgs_user Register_User

	// 解析 JSON 資料
	if err := c.ShouldBindJSON(&rgs_user); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid JSON format.",
		})
		return
	}

	password := rgs_user.Password
	checkPassword := rgs_user.CheckPassword

	//檢查使用者名稱是否已經被使用
	db := connectDB()
	var user User
	queryResult := db.Where("username = ?", rgs_user.Username).First(&user)
	if queryResult.Error == nil {
		c.JSON(400, gin.H{
			"message": "Username already exists.",
		})
		closeDB(db)
		return
	}

	// 檢查密碼強度
	if len(password) < 8 && len(password) < 20 {
		c.JSON(400, gin.H{
			"message": "Password should be at least 8 characters long.",
		})
		return
	}

	// 檢查兩次密碼是否相符
	if password != checkPassword {
		c.JSON(400, gin.H{
			"message": "Passwords do not match.",
		})
		return
	}

	// 檢查信箱格式
	if !isValidEmail(rgs_user.Email) {
		c.JSON(400, gin.H{
			"message": "Invalid email format.",
		})
		return
	}

	//檢查信箱是否已經被使用
	queryResult = db.Where("email = ?", rgs_user.Email).First(&user)
	if queryResult.Error == nil {
		c.JSON(400, gin.H{
			"message": "Email already exists.",
		})
		closeDB(db)
		return
	}

	closeDB(db)

	// 傳遞 rgs_user 給 CreateUser 函式
	CreateUser(c, rgs_user)
}

func HandleLogin(c *gin.Context) {
	var lg_user Login_User
	// 解析 JSON 資料
	if err := c.ShouldBindJSON(&lg_user); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid JSON format.",
		})
		return
	}

	log.Print(lg_user.Email)
	log.Print(lg_user.Password)
	db := connectDB()
	var user User
	queryResult := db.Where("email = ?", lg_user.Email).First(&user)
	if queryResult.Error != nil {
		c.JSON(400, gin.H{
			"message": "User does not exist.",
		})
		closeDB(db)
		return
	}

	if user.Password != lg_user.Password {
		c.JSON(400, gin.H{
			"message": "Incorrect password.",
		})
		closeDB(db)
		return
	}

	// 建立CustomClaims
	claims := CustomClaims{
		Username: lg_user.Email, // 使用者名稱
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(), // Token 有效期設定為 7 天
		},
	}

	// 產生Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	closeDB(db)
	c.JSON(200, gin.H{
		"message": "Login success.",
		"token":   signedToken,
		"userId":  user.Id,
	})
}

func HandleChangePassword(c *gin.Context) {
	var changePassword ChangePassword
	// 解析 JSON 資料
	if err := c.ShouldBindJSON(&changePassword); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid JSON format.",
		})
		return
	}

	db := connectDB()
	var user User
	queryResult := db.Where("id = ?", changePassword.UserID).First(&user)
	if queryResult.Error != nil {
		c.JSON(400, gin.H{
			"message": "User does not exist.",
		})
		closeDB(db)
		return
	}

	if user.Password != changePassword.OldPassword {
		c.JSON(400, gin.H{
			"message": "Incorrect password.",
		})
		closeDB(db)
		return
	}

	if changePassword.NewPassword != changePassword.ConfirmPassword {
		c.JSON(400, gin.H{
			"message": "Passwords do not match.",
		})
		closeDB(db)
		return
	}

	if len(changePassword.NewPassword) < 8 && len(changePassword.NewPassword) < 20 {
		c.JSON(400, gin.H{
			"message": "Password should be at least 8 characters long.",
		})
		closeDB(db)
		return
	}

	var newUser User
	newUser.Password = changePassword.NewPassword

	result := db.Model(&user).Where("id = ?", user.Id).Updates(&newUser)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "change password error" + result.Error.Error(),
		})
		closeDB(db)
		return
	}
	closeDB(db)
	c.JSON(200, gin.H{
		"message": "change password success",
	})
}

func CreateUser(c *gin.Context, rgs_user Register_User) {
	// 在這裡使用 rgs_user 中的資料進行相應的處理
	log.Print("create user")
	log.Print(rgs_user.Username)
	log.Print(rgs_user.Email)

	db := connectDB()
	var user User
	user.Id = uuid.New()
	user.Username = rgs_user.Username
	user.Password = rgs_user.Password
	user.Email = rgs_user.Email

	request := db.Create(&user)
	if request.Error != nil {
		c.JSON(500, gin.H{
			"message": "create user error" + request.Error.Error(),
		})
		closeDB(db)
		return
	}
	closeDB(db)

	c.JSON(200, gin.H{
		"message": "create user success",
	})
}

func UpdateUserById(c *gin.Context) {
	db := connectDB()
	var user User

	queryResult := db.Where("id = ?", c.Param("id")).First(&user)
	if queryResult.Error != nil {
		c.JSON(500, gin.H{
			"message": "query error" + queryResult.Error.Error(),
		})
		closeDB(db)
		return
	}

	var newUser User
	c.Bind(&newUser)
	newUser.Id = user.Id

	if !isValidEmail(newUser.Email) {
		c.JSON(400, gin.H{
			"message": "Invalid email format.",
		})
		return
	}

	result := db.Model(&user).Where("id = ?", user.Id).Updates(&newUser)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "update user error" + result.Error.Error(),
		})
		closeDB(db)
		return
	}
	closeDB(db)
	c.JSON(200, gin.H{
		"message": "update user success",
	})
}

func DeleteUserById(c *gin.Context) {
	db := connectDB()
	var user User

	queryResult := db.Where("id = ?", c.Param("id")).First(&user)
	if queryResult.Error != nil {
		c.JSON(500, gin.H{
			"message": "query error" + queryResult.Error.Error(),
		})
		closeDB(db)
		return
	}

	result := db.Delete(&user)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "delete user error" + result.Error.Error(),
		})
		closeDB(db)
		return
	}
	closeDB(db)
	c.JSON(200, gin.H{
		"message": "delete user success",
	})
}
