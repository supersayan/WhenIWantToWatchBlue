package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%q", html.EscapeString((r.URL.Path)))
}

func main() {
	server := NewServer()
	router := mux.NewRouter()

	router.Use(corsMiddleware)

	router.HandleFunc("/", home)
	router.HandleFunc("/tv/{index:[0-9]+}", server.getTvHandler).Methods("GET")
	router.HandleFunc("/tv/{index:[0-9]+}/flip", server.toggleTvHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

type Server struct {
	redis *redis.Client
	ctx   context.Context
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,emitempty"`
}

const REDIS_KEY = ""

func NewServer() *Server {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	return &Server{
		redis: rdb,
		ctx:   ctx,
	}
}

func (s *Server) getTvHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	indexStr := vars["index"]

	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 0 {
		s.sendResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid index parameter",
		})
		return
	}

	arrayStr, err := s.redis.Get(s.ctx, REDIS_KEY).Result()
	if err == redis.Nil {
		s.sendResponse(w, http.StatusOK, Response{
			Success: true,
			Data:    false,
		})
		return
	} else if err != nil {
		s.sendResponse(w, http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to retrieve data",
		})
		return
	}

	boolArray, err := s.parseArrayString(arrayStr)
	if err != nil {
		s.sendResponse(w, http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to parse data",
		})
		return
	}

	if index >= len(boolArray) {
		s.sendResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid index parameter",
		})
		return
	}

	s.sendResponse(w, http.StatusOK, Response{
		Success: true,
		Data:    boolArray[index],
	})
}

func (s *Server) toggleTvHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	indexStr := vars["index"]

	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 0 {
		s.sendResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid index parameter",
		})
		return
	}

	arrayStr, err := s.redis.Get(s.ctx, REDIS_KEY).Result()
	var boolArray []bool

	if err == redis.Nil {
		boolArray = make([]bool, 0)
	} else if err != nil {
		s.sendResponse(w, http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to retrieve data",
		})
		return
	} else {
		boolArray, err = s.parseArrayString(arrayStr)
		if err != nil {
			s.sendResponse(w, http.StatusInternalServerError, Response{
				Success: false,
				Error:   "Failed to parse data",
			})
			return
		}
	}

	if index >= len(boolArray) && index < 999 {
		newArray := make([]bool, index+1)
		copy(newArray, boolArray)
		boolArray = newArray
	}

	boolArray[index] = !boolArray[index]

	arrayStr = s.arrayToString(boolArray)
	err = s.redis.Set(s.ctx, REDIS_KEY, arrayStr, 0).Err()
	if err != nil {
		s.sendResponse(w, http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to save data",
		})
		return
	}

	s.sendResponse(w, http.StatusOK, Response{
		Success: true,
		Data:    boolArray[index],
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) sendResponse(w http.ResponseWriter, statusCode int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) arrayToString(arr []bool) string {
	strArr := make([]string, len(arr))
	for i, v := range arr {
		if v {
			strArr[i] = "1"
		} else {
			strArr[i] = "0"
		}
	}
	return strings.Join(strArr, "")
}

func (s *Server) parseArrayString(str string) ([]bool, error) {
	if str == "" {
		return []bool{}, nil
	}

	parts := strings.Split(str, "")
	boolArray := make([]bool, len(parts))

	for i, part := range parts {
		if part == "1" {
			boolArray[i] = true
		} else if part == "0" {
			boolArray[i] = false
		} else {
			return nil, fmt.Errorf("invalid boolean value: %s", part)
		}
	}

	return boolArray, nil
}
