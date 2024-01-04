package controller

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Comment struct {
	Id          uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Content     string    `json:"content"`
	User_id     uuid.UUID `gorm:"foreignKey:User_id" json:"user_id"`
	Article_id  uuid.UUID `gorm:"foreignKey:Article_id" json:"article_id"`
	PublishDate time.Time `json:"publish_date"`
}

type CommentwithAll struct {
	Comment
	User    *User    `json:"user,omitempty"`
	Article *Article `json:"article,omitempty"`
}

func (Comment) TableName() string {
	return "comment"
}

func GetComment(c *gin.Context) {
	db := connectDB()
	var comments []*CommentwithAll
	db.Preload("User").Preload("Article").Find(&comments)
	closeDB(db)
	c.JSON(200, comments)
}

func GetCommentById(c *gin.Context) {
	db := connectDB()
	var comment CommentwithAll
	queryResult := db.Preload("User").Preload("Article").Where("id = $1", c.Param("id")).Take(&comment)
	if queryResult.Error != nil {
		c.JSON(500, gin.H{
			"message": "query error" + queryResult.Error.Error(),
		})
		closeDB(db)
		return
	}
	closeDB(db)
	c.JSON(200, comment)
}

func GetCommentsByArticleId(c *gin.Context) {
	db := connectDB()
	var comments []*CommentwithAll
	queryResult := db.Preload("User").Preload("Article").Where("article_id = $1", c.Param("id")).Find(&comments)
	if queryResult.Error != nil {
		c.JSON(500, gin.H{
			"message": "query error" + queryResult.Error.Error(),
		})
		closeDB(db)
		return
	}
	closeDB(db)
	c.JSON(200, comments)
}

func CreateComment(c *gin.Context) {
	db := connectDB()
	var comment CommentwithAll
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(400, gin.H{
			"message": "bad request",
		})
		closeDB(db)
		return
	}
	comment.Id = uuid.New()
	comment.PublishDate = time.Now()
	queryResult := db.Create(&comment)
	if queryResult.Error != nil {
		c.JSON(500, gin.H{
			"message": "query error" + queryResult.Error.Error(),
		})
		closeDB(db)
		return
	}
	closeDB(db)
	c.JSON(200, gin.H{
		"message": "create comment success",
	})
}

func UpdateCommentById(c *gin.Context) {
	db := connectDB()
	var comment CommentwithAll
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(400, gin.H{
			"message": "bad request",
		})
		closeDB(db)
		return
	}
	comment.PublishDate = time.Now()
	queryResult := db.Save(&comment)
	if queryResult.Error != nil {
		c.JSON(500, gin.H{
			"message": "update error" + queryResult.Error.Error(),
		})
		closeDB(db)
		return
	}
	closeDB(db)
	c.JSON(200, gin.H{
		"message": "update comment success",
	})
}

func DeleteCommentById(c *gin.Context) {
	db := connectDB()
	var comment Comment
	queryResult := db.Where("id = ?", c.Param("id")).Delete(&comment)
	if queryResult.Error != nil {
		c.JSON(500, gin.H{
			"message": "delete error" + queryResult.Error.Error(),
		})
		closeDB(db)
		return
	}
	closeDB(db)
	c.JSON(200, gin.H{
		"message": "delete comment success",
	})
}
