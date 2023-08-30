package user

type User struct {
	UID      int    `json:"uid" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	License  string `json:"license"`
}
