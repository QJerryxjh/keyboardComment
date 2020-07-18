package managerController

import (
	"keyboardComment/dbs"
	"keyboardComment/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetManagers(c *gin.Context) {
	var managers []model.DbManager
	err := dbs.DB.Find(&managers).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "failed",
			"err":  err,
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
		"code":     200,
		"msg":      "success",
		"managers": resManagers,
	})
}
