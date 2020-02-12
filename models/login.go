package models

type User struct {
  Password string `json:"password"`
}

func (u *User) UserCheck(password string) (bool, error) {
  return true, nil
}