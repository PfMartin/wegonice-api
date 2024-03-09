package db

import "go.mongodb.org/mongo-driver/mongo/options"

type Pagination struct {
	PageID   int64 `form:"page_id" json:"page_id" binding:"min=1"`
	PageSize int64 `form:"page_id" json:"page_size" binding:"min=0,max=500"`
}

func (pagination *Pagination) getFindOptions() *options.FindOptions {
	findOptions := options.Find()
	findOptions.SetSkip((pagination.PageID - 1) * pagination.PageSize)
	findOptions.SetLimit(pagination.PageSize)

	return findOptions
}

type Role string

const (
	UserRole  Role = "user"
	AdminRole Role = "admin"
)

type User struct {
	ID         string `json:"id" bson:"_id"`
	Email      string `json:"email" bson:"email"`
	Password   string `json:"password" bson:"password"`
	Role       Role   `json:"role" bson:"role"`
	CreatedAt  int64  `json:"createdAt" bson:"createdAt"`
	ModifiedAt int64  `json:"modifiedAt" bson:"modifiedAt"`
}
