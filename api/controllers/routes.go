package controllers

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) initializeRoutes() {

	v1 := s.Router.Group("/api/v1")
	{
		// Inverter routes
		v1.GET("/inverters", s.GetInverters)
		v1.GET("/inverters/:serial", s.GetInverter)
	}
}

func answerJSON(c *gin.Context, code int, payload interface{}) {
	c.JSON(code, payload)
}

func answerError(c *gin.Context, code int, message string) {
	answerJSON(c, code, map[string]string{"error": message})
}
