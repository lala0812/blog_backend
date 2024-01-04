package controller

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type Article struct {
	Id          uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	User_id     uuid.UUID `gorm:"foreignKey:User_id" json:"user_id"`
	PublishDate time.Time `json:"publish_date"`
}

func (Article) TableName() string {
	return "articles"
}

type ArticlewithAll struct {
	Article
	User *User `json:"user,omitempty"`
}

func GetArticle(c *gin.Context) {
	db := connectDB()
	var articles []*ArticlewithAll
	db.Preload(clause.Associations).Find(&articles)
	closeDB(db)
	c.JSON(200, articles)
}

func GetArticleById(c *gin.Context) {
	db := connectDB()
	var article ArticlewithAll
	queryResult := db.Preload(clause.Associations).Where("id = $1", c.Param("id")).Take(&article)
	if queryResult.Error != nil {
		c.JSON(500, gin.H{
			"message": "query error" + queryResult.Error.Error(),
		})
		closeDB(db)
		return
	}
	closeDB(db)
	c.JSON(200, article)
}

func GetArticleByUserId(c *gin.Context) {
	db := connectDB()
	var articles []*ArticlewithAll
	queryResult := db.Preload(clause.Associations).Where("user_id = $1", c.Param("id")).Find(&articles)
	if queryResult.Error != nil {
		c.JSON(500, gin.H{
			"message": "query error" + queryResult.Error.Error(),
		})
		closeDB(db)
		return
	}
	closeDB(db)
	c.JSON(200, articles)
}

func CreateArticle(c *gin.Context) {
	db := connectDB()
	var article Article
	c.BindJSON(&article)
	article.Id = uuid.New()
	article.PublishDate = time.Now()
	request := db.Create(&article)
	if request.Error != nil {
		c.JSON(500, gin.H{
			"message": "create article error" + request.Error.Error(),
		})
		closeDB(db)
		return
	}
	closeDB(db)
	c.JSON(200, gin.H{
		"message": "create article success",
	})
}

func UpdateArticleById(c *gin.Context) {
	db := connectDB()
	var article Article
	queryResult := db.Where("id = ?", c.Param("id")).First(&article)
	if queryResult.Error != nil {
		c.JSON(500, gin.H{
			"message": "query error" + queryResult.Error.Error(),
		})
		closeDB(db)
		return
	}
	c.BindJSON(&article)
	request := db.Save(&article)
	if request.Error != nil {
		c.JSON(500, gin.H{
			"message": "update error" + request.Error.Error(),
		})
		closeDB(db)
		return
	}
	closeDB(db)
	c.JSON(200, gin.H{
		"message": "update article success",
	})
}

func DeleteArticleById(c *gin.Context) {
	db := connectDB()
	var article Article
	queryResult := db.Where("id = ?", c.Param("id")).First(&article)
	if queryResult.Error != nil {
		c.JSON(500, gin.H{
			"message": "query error" + queryResult.Error.Error(),
		})
		closeDB(db)
		return
	}
	request := db.Delete(&article)
	if request.Error != nil {
		c.JSON(500, gin.H{
			"message": "delete error" + request.Error.Error(),
		})
		closeDB(db)
		return
	}
	closeDB(db)
	c.JSON(200, gin.H{
		"message": "delete success",
	})
}
