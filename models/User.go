package models

type User struct {
	Email        string `json:"email" bson:"email"`
	Nickname     string `json:"nickname" bson:"nickname"`
	ProfileImage string `json:"profile_image" bson:"profile_image"`
	Password     string `json:"password" bson:"password"`
	Type         string `json:"type" bson:"type"`
	CreatedAt    string `json:"created_at" bson:"created_at"`
	AccessToken  string `json:"access_token" bson:"access_token"`
}
