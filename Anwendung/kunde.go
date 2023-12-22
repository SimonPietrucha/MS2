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

type Kunde struct {
	Collection *mongo.Collection // Füge die Collection hinzu
}

type KundeModel struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Quantity int                `json:"quantity,omitempty" bson:"quantity,omitempty"`
}

func (p *Kunde) Create(w http.ResponseWriter, r *http.Request) {
	var kunde KundeModel
	if err := json.NewDecoder(r.Body).Decode(&kunde); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := p.Collection.InsertOne(context.TODO(), kunde)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result.InsertedID)
}

func (p *Kunde) List(w http.ResponseWriter, r *http.Request) {
	cursor, err := p.Collection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var kunde []KundeModel
	if err := cursor.All(context.TODO(), &kunde); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(kunde)
}

func (p *Kunde) GetByID(w http.ResponseWriter, r *http.Request) {
	kundeID := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(kundeID)
	if err != nil {
		http.Error(w, "Ungültige Kunden ID", http.StatusBadRequest)
		return
	}

	var kunde KundeModel
	err = p.Collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&kunde)
	if err != nil {
		http.Error(w, "Kunde konnte nicht gefunden werden", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(kunde)
}

func (p *Kunde) UpdateByID(w http.ResponseWriter, r *http.Request) {
	kundeID := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(kundeID)
	if err != nil {
		http.Error(w, "Ungültige Kunden ID", http.StatusBadRequest)
		return
	}

	var updatedKunde KundeModel
	if err := json.NewDecoder(r.Body).Decode(&updatedKunde); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = p.Collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objID},
		bson.M{"$set": updatedKunde},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (p *Kunde) DeleteByID(w http.ResponseWriter, r *http.Request) {
	kundeID := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(kundeID)
	if err != nil {
		http.Error(w, "Ungültige Kunden ID", http.StatusBadRequest)
		return
	}

	_, err = p.Collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
