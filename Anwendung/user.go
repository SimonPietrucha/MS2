package Anwendung

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Collection *mongo.Collection // Füge die Collection hinzu
}

type UserModel struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Quantity int                `json:"quantity,omitempty" bson:"quantity,omitempty"`
}

func (p *User) Create(w http.ResponseWriter, r *http.Request) {
	var user UserModel
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := p.Collection.InsertOne(context.TODO(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result.InsertedID)
}

func (p *User) List(w http.ResponseWriter, r *http.Request) {
	cursor, err := p.Collection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var user []UserModel
	if err := cursor.All(context.TODO(), &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (p *User) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user UserModel
	err = p.Collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (p *User) UpdateByID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var updatedUser UserModel
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = p.Collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objID},
		bson.M{"$set": updatedUser},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (p *User) DeleteByID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	_, err = p.Collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
