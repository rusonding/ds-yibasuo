package models

type Response struct {
  Code    int64       `json:"code"`
  Message string      `json:"message"`
  Result  interface{} `json:"result"`
}