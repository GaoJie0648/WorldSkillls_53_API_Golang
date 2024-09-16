package utils

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Create(client *mongo.Client, dbName string, collectionName string, data interface{}) interface{} {
	collection := client.Database(dbName).Collection(collectionName)
	response, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		return err
	}
	return response.InsertedID
}

func Read(client *mongo.Client, dbName string, collectionName string, filter bson.M) map[string]interface{} {
	collection := client.Database(dbName).Collection(collectionName)
	var result map[string]interface{}
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil
	}
	return result
}

func ReadAll(client *mongo.Client, dbName string, collectionName string, filter bson.M) []map[string]interface{} {
	collection := client.Database(dbName).Collection(collectionName)
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil
	}
	var results []map[string]interface{}
	for cursor.Next(context.Background()) {
		var result map[string]interface{}
		err := cursor.Decode(&result)
		if err != nil {
			return nil
		}
		results = append(results, result)
	}
	return results
}

func Update(client *mongo.Client, dbName string, collectionName string, filter bson.M, update bson.M) error {
	collection := client.Database(dbName).Collection(collectionName)
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func Delete(client *mongo.Client, dbName string, collectionName string, filter bson.M) error {
	collection := client.Database(dbName).Collection(collectionName)
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}
