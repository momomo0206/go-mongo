package services

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Todo struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	Task      string    `json:"task,omitempty" bson:"task,omitempty"`
	Completed bool      `json:"completed" bson:"completed"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

var client *mongo.Client

func New(mongo *mongo.Client) Todo {
	client = mongo
	return Todo{}
}

func returnCollectionPointer(collection string) *mongo.Collection {
	return client.Database("todo").Collection(collection)
}

func (t *Todo) GetAllTodos() ([]Todo, error) {
	collection := returnCollectionPointer("todos")
	var todos []Todo

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo
		cursor.Decode(&todo)
		todos = append(todos, todo)
	}

	return todos, nil
}

func (t *Todo) GetTodoById(id string) (Todo, error) {
	collection := returnCollectionPointer("todos")
	var todo Todo

	mongoID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Todo{}, err
	}

	err = collection.FindOne(context.Background(), bson.M{"_id": mongoID}).Decode(&todo)
	if err != nil {
		log.Println(err)
		return Todo{}, err
	}

	return todo, nil
}

func (t *Todo) InsertTodo(entry Todo) error {
	collection := returnCollectionPointer("todos")

	_, err := collection.InsertOne(context.TODO(), Todo{
		Task:      entry.Task,
		Completed: entry.Completed,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error:", err)
		return err
	}

	return nil
}

func (t *Todo) UpdateTodo(id string, entry Todo) (*mongo.UpdateResult, error) {
	collection := returnCollectionPointer("todos")
	mongoID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	update := bson.D{
		{Key: "set", Value: bson.D{
			{Key: "task", Value: entry.Task},
			{Key: "completed", Value: entry.Completed},
			{Key: "updated_at", Value: time.Now()},
		}},
	}

	res, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": mongoID},
		update,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return res, nil
}

func (t *Todo) DeleteTodo(id string) error {
	collection := returnCollectionPointer("todos")
	mongoID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = collection.DeleteOne(
		context.Background(),
		bson.M{"_id": mongoID},
	)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
