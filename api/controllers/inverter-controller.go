package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rjmalves/cpid-solar-gateway/api/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetInverters : handler for listing all the inverters in the DB
func (s *Server) GetInverters(c *gin.Context) {
	inverters, err := models.ListInverters(s.DB)
	if err != nil {
		answerError(c, http.StatusInternalServerError, err.Error())
		return
	}
	answerJSON(c, http.StatusOK, inverters)
}

// GetInverter : handler for getting a single inverter from the DB
func (s *Server) GetInverter(c *gin.Context) {
	serial := c.Param("serial")
	i := &models.Inverter{
		Serial: serial,
	}
	if err := i.ReadInverter(s.DB); err != nil {
		if err == mongo.ErrNoDocuments {
			answerError(c, http.StatusNotFound, err.Error())
		} else {
			answerError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	answerJSON(c, http.StatusOK, i)
}
