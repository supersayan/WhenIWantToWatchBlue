package main

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

const REDIS_TV_KEY = "tv"

func (s *Server) initRedisTvs() error {
	_, err := s.redis.Set(s.ctx, REDIS_TV_KEY, "\x00\x00", 0).Result()
	return err
}

func (s *Server) getRedisTvsLength() (int, error) {
	length, err := s.redis.StrLen(s.ctx, REDIS_TV_KEY).Result()
	if err != nil {
		return 0, err
	}
	length *= 8
	return int(length), nil
}

func (s *Server) getRedisTvs(initIfNil bool) ([]bool, error) {
	arrayStr, err := s.redis.Get(s.ctx, REDIS_TV_KEY).Result()
	var tvArray []bool

	if err == redis.Nil {
		if initIfNil {
			err = s.initRedisTvs()
			if err != nil {
				return nil, errors.New("Failed to initialize data")
			}
			return s.getRedisTvs(false)
		} else {
			return nil, errors.New("Data not initialized")
		}
	} else if err != nil {
		return nil, errors.New("Failed to retrieve data")
	} else {
		tvArray, err = s.parseArrayString(arrayStr)
		if err != nil {
			return nil, errors.New("Failed to parse data")
		}
	}

	return tvArray, nil
}

func (s *Server) getRedisTv(index int, initIfNil bool) (bool, error) {
	tvArrayLen, err := s.getRedisTvsLength()
	if err != nil {
		return false, err
	}
	if index > tvArrayLen {
		return false, errors.New("Index out of bounds")
	}

	tvArray, err := s.getRedisTvs(initIfNil)

	if err != nil {
		return false, err
	}

	tv := tvArray[index]

	return tv, nil
}

func (s *Server) setRedisTv(index int, value bool) (bool, error) {
	tvArrayLen, err := s.getRedisTvsLength()

	if err != nil {
		return false, err
	}
	if index > tvArrayLen {
		return false, errors.New("Index out of bounds")
	}

	var valueInt int
	if value {
		valueInt = 1
	} else {
		valueInt = 0
	}
	_, err = s.redis.SetBit(s.ctx, REDIS_TV_KEY, int64(index), valueInt).Result()
	if err != nil {
		return false, errors.New("Failed to set data")
	}

	return value, nil
}

func (s *Server) toggleRedisTv(index int) (bool, error) {
	tv, err := s.getRedisTv(index, true)

	if err != nil {
		return false, err
	}

	return s.setRedisTv(index, !tv)
}
