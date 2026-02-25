package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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

	res, err := s.getRedisTv(index, false)
	if err != nil {
		s.sendResponse(w, http.StatusInternalServerError, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	s.sendResponse(w, http.StatusOK, Response{
		Success: true,
		Data:    res,
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

	res, err := s.toggleRedisTv(index)
	if err != nil {
		s.sendResponse(w, http.StatusInternalServerError, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	s.sendResponse(w, http.StatusOK, Response{
		Success: true,
		Data:    res,
	})
}

func (s *Server) getAllTvHandler(w http.ResponseWriter, r *http.Request) {
	arrayStr, err := s.getRedisTvs(true)
	if err != nil {
		s.sendResponse(w, http.StatusInternalServerError, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	s.sendResponse(w, http.StatusOK, Response{
		Success: true,
		Data:    arrayStr,
	})
}
