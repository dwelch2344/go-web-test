package main

import (
	"github.com/derekdowling/go-json-spec-handler"
	"golang.org/x/net/context"
	"github.com/derekdowling/go-json-spec-handler/jsh-api/store"
	"log"
	"reflect"
	"container/list"
)

type User struct {
	ID string
	Name string `json:"name"`
}

type UserStorage struct {
	store.CRUD
}


func convert(list *list.List) (*jsh.List, error){
	return nil, nil
}

func (us *UserStorage) Get(ctx context.Context, id string) (*jsh.Object, jsh.ErrorType){
	user := &User{
		ID: "123",
		Name: "Dave",
	}


	object, _ := jsh.NewObject(user.ID, reflect.TypeOf(*user).Name(), user)
	return object, nil
}
func (us *UserStorage) List(ctx context.Context) (jsh.List, jsh.ErrorType){
	log.Println("Test")
	user := &User{
		ID: "123",
		Name: "Dave",
	}


	object, _ := jsh.NewObject(user.ID, reflect.TypeOf(*user).Name(), user)
	//err := object.Unmarshal("user", user)
	//if err != nil {
	//	log.Printf("Hello world %v \n", err)
	//}

	testList := jsh.List{object}

	return testList, nil
}

func (us *UserStorage) Delete(ctx context.Context, id string) jsh.ErrorType{
	return nil
}

func (us *UserStorage) Save(ctx context.Context, object *jsh.Object) (*jsh.Object, jsh.ErrorType) {
	user := &User{}
	err := object.Unmarshal("user", user)
	if err != nil {
		return nil, err
	}

	// generate your id, however you choose
	user.ID = "1234"

	//err := object.Marshal(user)
	//if err != nil {
	//	return nil, err
	//}

	// do save logic
	return object, nil
}

func (us *UserStorage) Update(ctx context.Context, object *jsh.Object) (*jsh.Object, jsh.ErrorType) {
	user := &User{}
	err := object.Unmarshal("user", user)
	if err != nil {
		return nil, err
	}

	user.Name = "NewName"

	//err := object.Marshal(user)
	//if err != nil {
	//	return nil, err
	//}

	// perform patch
	return object, nil
}