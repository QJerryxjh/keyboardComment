package managerController

import (
	"keyboardComment/dbs"
	"keyboardComment/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 获取管理员列表
func GetManagers(c *gin.Context) {
	var managers []model.DbManager
	err := dbs.DB.Find(&managers).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "failed",
			"err":  err.Error(),
		})
		return
	}

	resManagers := make([]map[string]interface{}, 0, len(managers))
	for _, value := range managers {
		resManager := map[string]interface{}{
			"id":         value.ID,
			"name":       value.Name,
			"account":    value.Account,
			"created_at": value.CreatedAt.Format("2006-01-02"),
			"gender":     value.Gender,
		}
		resManagers = append(resManagers, resManager)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": resManagers,
	})
}

// 根据id获取管理员信息
func GetManagerInfo(c *gin.Context) {
	managerId := c.PostForm("id")
	if managerId == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "缺少id字段",
		})
		return
	}
	var manager model.DbManager
	err := dbs.DB.Where("id = ?", managerId).First(&manager).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "failed",
			"err":  err.Error(),
		})
		return
	}
	resManager := map[string]interface{}{
		"id":         manager.ID,
		"name":       manager.Name,
		"account":    manager.Account,
		"created_at": manager.CreatedAt.Format("2006-01-02"),
		"gender":     manager.Gender,
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": resManager,
	})
}

// 根据id删除管理员条目
func DeleteManager(c *gin.Context) {
	managerId := c.PostForm("id")
	if managerId == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "缺少id字段",
		})
		return
	}
	var manager model.DbManager
	err := dbs.DB.Where("id = ?", managerId).First(&manager).Delete(&manager).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "failed",
			"err":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
	})
}
