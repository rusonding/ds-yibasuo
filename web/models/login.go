package models

import (
  "encoding/json"
)

type User struct {
  Password string `json:"password"`
}

func (u *User) UserCheck(request []byte) (*Response, error) {
  if err := json.Unmarshal(request, u); err != nil {
    return &Response{Code:500, Message:"参数错误",Result:nil}, nil
  } else {
    return nil, err
  }
}