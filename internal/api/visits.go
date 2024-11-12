package api

import (
	"fmt"
	"math"
	"net/http"

	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type BasicVisitsResponse struct {
	visits  []models.Visit
	current int
}

func GetTrafficData() echo.HandlerFunc {
	return func(c echo.Context) error {
		visits, err := models.GetVisits()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching visits: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, BasicVisitsResponse{visits: visits, current: models.CurrentVisitors()})
	}
}

func GetVisits() echo.HandlerFunc {
	return func(c echo.Context) error {
		qualityStr := c.QueryParam("quality")
		timeframeStr := c.QueryParam("timeframe")

		visits, err := models.GetVisits()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching visits: %v", err), Errors: []string{err.Error()}})
		}

		if len(qualityStr) < 1 || len(timeframeStr) < 1 {
			return c.JSON(http.StatusOK, BasicVisitsResponse{visits: visits, current: models.CurrentVisitors()})
		}

		quality := models.ParseVisitQuality(qualityStr)

		timeframe := models.ParseTimeframe(timeframeStr)

		data, err := models.GetVisitsByQualityAndTimeframe(quality, timeframe)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching visits data: %v", err), Errors: []string{err.Error()}})
		}

		totVis := len(visits)

		unique, err := models.CountUniqueIps()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error counting unique visitors: %v", err), Errors: []string{err.Error()}})
		}

		views, err := models.CountTotalViews()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error counting total views: %v", err), Errors: []string{err.Error()}})
		}

		avgDuration, err := models.GetAverageVisitDuration()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error counting average visit duration: %v", err), Errors: []string{err.Error()}})
		}

		zeroers, err := models.GetVisitsWithZeroViews()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error counting visits with zero views: %v", err), Errors: []string{err.Error()}})
		}

		vsOrigins, err := models.GetVisitOrigins()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error counting visit origins: %v", err), Errors: []string{err.Error()}})
		}

		dvOrigins, err := models.GetDeviceOrigins()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error counting device origins: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, models.VisitsResponse{Data: data, Current: models.CurrentVisitors(), TotalUniqueVisitors: unique, TotalViews: views, BounceRate: fmt.Sprintf("%d%%", int(math.Floor(float64(zeroers)/float64(totVis)*100))), AvgVisitDuration: avgDuration, TotalVisits: totVis, VisitOrigins: vsOrigins, DeviceOrigins: dvOrigins})
	}
}

func GetVisitStats() echo.HandlerFunc {
	return func(c echo.Context) error {
		visits, err := models.GetVisits()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching visits: %v", err), Errors: []string{err.Error()}})
		}
		totVis := len(visits)

		unique, err := models.CountUniqueIps()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error counting unique visitors: %v", err), Errors: []string{err.Error()}})
		}

		views, err := models.CountTotalViews()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error counting total views: %v", err), Errors: []string{err.Error()}})
		}

		avgDuration, err := models.GetAverageVisitDuration()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error counting average visit duration: %v", err), Errors: []string{err.Error()}})
		}

		zeroers, err := models.GetVisitsWithZeroViews()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error counting visits with zero views: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, models.VisitsStats{Current: models.CurrentVisitors(), TotalViews: views, BounceRate: fmt.Sprintf("%d%%", int(math.Floor(float64(zeroers)/float64(totVis)*100))), AvgVisitDuration: avgDuration, TotalVisits: totVis, TotalUniqueVisitors: unique})
	}
}

func GetVisitsStandings() echo.HandlerFunc {
	return func(c echo.Context) error {
		vsOrigins, err := models.GetVisitOrigins()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error counting visit origins: %v", err), Errors: []string{err.Error()}})
		}

		dvOrigins, err := models.GetDeviceOrigins()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error counting device origins: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusOK, models.VisitsStandings{VisitOrigins: vsOrigins, DeviceOrigins: dvOrigins})
	}
}

func GetVisitGraph() echo.HandlerFunc {
	return func(c echo.Context) error {
		qualityStr := c.QueryParam("quality")
		timeframeStr := c.QueryParam("timeframe")

		quality := models.ParseVisitQuality(qualityStr)

		timeframe := models.ParseTimeframe(timeframeStr)

		data, err := models.GetVisitsByQualityAndTimeframe(quality, timeframe)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.JSONErrorResponse{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error fetching visits data: %v", err), Errors: []string{err.Error()}})
		}

		log.Infof("Visits data: %v", data)

		return c.JSON(http.StatusOK, models.Graph{Data: data})
	}
}
