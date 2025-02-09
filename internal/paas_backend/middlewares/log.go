package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Log(c *gin.Context) {
	c.Next()
	logrus.Infof("Request: %v %v", c.Writer.Status(), c.Request.URL)

}
