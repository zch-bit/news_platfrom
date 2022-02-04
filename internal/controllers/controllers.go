package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"newsplatform/internal/database"
	"newsplatform/internal/models"
)

type InputNews struct {
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	PublishedAt time.Time `json:"publishedAt"`
	Content     string    `json:"content"`
}

// TODO support pagination
func FindNews(c *gin.Context) {
	var news []models.News
	database.DB.Find(&news)
	c.JSON(http.StatusOK, gin.H{
		"news": news,
	})
}

func CreateNews(c *gin.Context) {
	var input InputNews
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	article := models.News{
		Title:       input.Title,
		Author:      input.Author,
		Description: input.Description,
		URL:         input.URL,
		PublishedAt: input.PublishedAt,
		Content:     input.Content,
	}
	database.DB.Create(&article)
	c.JSON(http.StatusOK, gin.H{
		"data": input,
	})
}

// TODO support pagination
func FindNew(c *gin.Context) {
	var news []models.News
	if err := database.DB.Where("title like '%%%s'", c.Param("keyword")).Find(&news).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "no such record",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"news": news,
	})
}

func UpdateNew(c *gin.Context) {
	var n models.News
	if err := database.DB.Where("id = ?", c.Param("id")).First(&n).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "no such record",
		})
		return
	}

	var input InputNews
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	database.DB.Model(&n).Updates(input)
	c.JSON(http.StatusOK, gin.H{
		"news": n,
	})
}

func DeleteNews(c *gin.Context) {
	var news []models.News

	id := c.Param("id")
	if err := database.DB.Where("id = ?", id).Find(&news).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "no such record",
		})
	}

	database.DB.Delete(&news)
	c.JSON(http.StatusOK, gin.H{
		"news": news,
	})
}

func DeleteNewsByDay(c *gin.Context) {
	var news []models.News
	n, err := strconv.ParseInt(c.Param("day"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}
	beforeDate := time.Now().Add(time.Duration(-1*n*24) * time.Hour)
	if err := database.DB.Where("published_at < ?", beforeDate).Find(&news).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	database.DB.Delete(&news)
	c.JSON(http.StatusOK, gin.H{
		"news": news,
	})
}
