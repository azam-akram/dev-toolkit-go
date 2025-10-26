package dynamo_db

import "github/dev-toolkit-go/aws-go/aws-apigateway-lambda-dynamo-go/model"

type DBHandler interface {
	Save(book *model.MyBook) error
	Update(book *model.MyBook) error
	UpdateAttributeByID(id, key, value string) error
	GetByID(id string) (*model.MyBook, error)
	DeleteByID(id string) error
}
