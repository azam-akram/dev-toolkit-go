package utils_test

import (
	"dev-toolkit-go/docker-k8s-kafka-mongodb-go/internal/model"
	"dev-toolkit-go/docker-k8s-kafka-mongodb-go/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringToStruct(t *testing.T) {
	assertThat := assert.New(t)

	str := `{
		"id": 1,
		"name": "The User"
	}`

	var student model.Student
	err := utils.StringToStruct(str, &student)

	assertThat.Nil(err)
	assertThat.Equal(1, student.ID)
	assertThat.Equal("The User", student.Name)
}
