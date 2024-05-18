package dynamo_db

import "github/dev-toolkit-go/aws-go/apigateway-lambda-dynamo-aws-go/model"

type Handler interface {
	Save(book *model.MyBook) error
	Update(book *model.MyBook) error
	UpdateAttributeByID(id, key, value string) error
	GetByID(id string) (*model.MyBook, error)
	DeleteByID(id string) error
}
