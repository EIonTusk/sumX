// Package api provides the HTTP API for the summarization service
package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"sumx/db"
	"sumx/llm"
	"sumx/x"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
)

type Config struct {
	XBearerToken string
	frontendBasePath string
}

type Server struct {
	Router 	*gin.Engine
	LLM    	*llm.HFClient
	DB		*bun.DB
	Config  *Config
}

type SummaryFrontend struct{
	Params SummaryParams `json:"params"`
	Summary []db.SummaryData `json:"summary"`
	Tweets []string			`json:"tweets"`
}

type SummaryParams struct{
	Username string `json:"username"`
	From string 	`json:"from"`
	To string 		`json:"to"`
	Limit int64 	`json:"limit"`
}


// NewServer creates a new API server instance
func NewServer(llmClient *llm.HFClient, db *bun.DB) *Server {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	frontendBasePath := os.Getenv("FRONTEND_BASE_PATH")
	if frontendBasePath == "" {
		panic("FRONTEND_BASE_PATH is empty")
	}

	bearerToken := os.Getenv("X_BEARER_TOKEN")
	if bearerToken == "" {
		panic("X_BEARER_TOKEN is empty")
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
			AllowOrigins:     []string{frontendBasePath}, // Frontend dev server
			AllowMethods:     []string{"GET", "POST", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	s := &Server{
		Router: r,
		LLM:    llmClient,
		DB: 	db,
		Config: &Config{
			XBearerToken: bearerToken,
			frontendBasePath: frontendBasePath,
		},
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.Router.GET("/api/summarize/:username", s.handleSummarize)
	s.Router.GET("/api/summaries", s.handleGetSummaries)
}

func (s *Server) handleGetSummaries(c *gin.Context) {
	var data []db.Summary
	err := s.DB.NewSelect().Model(&data).Relation("XUser").Order("created_at DESC").Scan(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Summaries could not be loaded"})
		return
	}
	var formatedData []SummaryFrontend
	for _, d := range data {
		from := ""
		to := ""
		if !d.From.IsZero() {
			from = d.From.Format(time.RFC3339)
		} 
		if !d.To.IsZero() {
			to = d.To.Format(time.RFC3339)
		}
		formatedData = append(formatedData, SummaryFrontend{
			Params: SummaryParams{
				Username: d.XUser.UserName,
				From: from,
				To: to,
				Limit: int64(d.Limit),
			},
			Summary: d.Summary,
			Tweets: d.Tweets,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": formatedData,
	})
}

func (s *Server) handleSummarize(c *gin.Context) {
	username := c.Param("username")
	if username[0] == '@' {
		username = username[1:]
	}
	limitAsString := c.DefaultQuery("limit", "-1")
	limit, err := strconv.Atoi(limitAsString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Limit has to be a integer"})
		return
	}
	if limit != -1 && (limit < 5 || limit > 100) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Limit not in range [5, 100]"})
		return
	}

	from := c.DefaultQuery("from", "")
	to := c.DefaultQuery("to", "")

	var xUser db.XUser
	var userID string

	err = s.DB.NewSelect().Model(&xUser).Where("user_name = ?", username).Scan(context.Background())
	if errors.Is(err, sql.ErrNoRows) {
		userID, err = x.FetchUserID(username, s.Config.XBearerToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}

		userIDInt64, err := strconv.ParseInt(userID, 10, 64)
		if err == nil {
			_ = s.DB.NewInsert().Model(&db.XUser{
				ID: userIDInt64,
				UserName: username,
			}).Scan(context.Background())
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query db"})
		return
	} else {
		userID = strconv.FormatInt(xUser.ID, 10)
	}

	tweets, nextReset, err := x.FetchTweetsByUsernameTimeframe(userID, from, to, limit, s.Config.XBearerToken)
	if err != nil {
		nextResetInt, err1 := strconv.ParseInt(nextReset, 10, 64)
		if err1 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error()})
			return
		}
		rTime := time.Unix(nextResetInt, 0).UTC().Format(time.RFC3339)
		c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error(), "next_reset": rTime})
		return
	}

	rawTweets := tweetsToText(tweets)
	data, err := s.LLM.Summarize(rawTweets)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to summarize tweets"})
		return
	}

	var jsonData []db.SummaryData
	err = json.Unmarshal( []byte(data), &jsonData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to summarize tweets"})
		return 
	}

	tweetsAsString := make([]string, len(tweets))
	for i, t := range tweets {
		tweetsAsString[i] = t.Text
	} 

	userIDInt64, err1 := strconv.ParseInt(userID, 10, 64)
	fromTime, err2 := time.Parse(time.RFC3339, from)
	toTime, err3 := time.Parse(time.RFC3339, to)
	if err1 != nil || err2 != nil || err3 != nil {
		_ = s.DB.NewInsert().Model(&db.Summary{
			XUserID: userIDInt64,
			From: fromTime,
			To: toTime,
			Limit: int16(limit),
			Summary: jsonData,
			Tweets: tweetsAsString,
		}).Scan(context.Background())
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"params": gin.H{
				"username": username, 
				"from": from,
				"to": to,
				"limit": limit, 
			},
			"summary": jsonData,
			"tweets": tweetsAsString,
		},
	})
}

func tweetsToText(tweets []x.Tweet) string {
	text := ""
	for i, tweet := range tweets {
if i == len(tweets) -1 {
			text += tweet.Text
			break
		}
		text += tweet.Text + "\n---\n"
	}
	return text
}
