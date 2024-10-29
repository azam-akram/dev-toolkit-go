package dynamo_db

import "github.com/dev-toolkit-go/aws-apigateway-lambda-demo-go/model"

type DBHandler interface {
	Save(book *model.MyBook) error
	Update(book *model.MyBook) error
	UpdateAttributeByID(id, key, value string) error
	GetByID(id string) (*model.MyBook, error)
	DeleteByID(id string) error
}
