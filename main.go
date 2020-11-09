package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
)

type Article struct{
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	TITLE string `json:"title,omitempty" bson:"title,omitempty"`
	SUBTITLE string `json:"subtitle,omitempty" bson:"subtitle,omitempty"`
	CONTENT string `json:"content,omitempty" bson:"content,omitempty"`
	TIME  time.Time `json:"time,omitempty" bson:"time,omitempty"`

}
var articles []Article
var client *mongo.Client

func createArticle(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	var article Article
	_=json.NewDecoder(r.Body).Decode(&article)
	collection :=client.Database("inshorts").Collection("article")
	ctx, _ := context.WithTimeout(context.Background(),  10*time.Second)
	result, _ :=collection.InsertOne(ctx,article)
	//article.ID=strconv.Itoa(rand.Intn(100000)) //generate random id
	//articles=append(articles, article)
	json.NewEncoder(w).Encode(result)



}
func getAll(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	var articles []Article
	collection :=client.Database("inshorts").Collection("article")
	ctx, _ := context.WithTimeout(context.Background(),  10*time.Second)
	cursor, err :=collection.Find(ctx, bson.M{})
	if err !=nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error()+`"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx){
		var article Article
		cursor.Decode(&article)
		articles=append(articles,article)
	}
	if err :=cursor.Err(); err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error()+`"}`))
		return
	}
	json.NewEncoder(w).Encode(articles)


}
func getArticleByID(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	params :=mux.Vars(r)
	id, _ :=primitive.ObjectIDFromHex(params["id"])
	var articles Article
	collection :=client.Database("inshorts").Collection("article")
	ctx, _ := context.WithTimeout(context.Background(),  10*time.Second)
	err := collection.FindOne(ctx, Article{ID: id}).Decode(&articles)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error()+`"}`))
		return
	}
	json.NewEncoder(w).Encode(articles)


}
func main(){
	fmt.Println("Hello World")
	ctx, _ := context.WithTimeout(context.Background(),  10*time.Second)
	client, _ = mongo.Connect(ctx,options.Client().ApplyURI("mongodb://localhost:27017"))
	router:=mux.NewRouter()
	//articles=append(articles, Article{ID:"1", TITLE:"Hello",SUBTITLE:"None",CONTENT:"You know"})

	router.HandleFunc("/articles",createArticle).Methods("POST")
	router.HandleFunc("/articles/{id}", getArticleByID).Methods("GET")
	router.HandleFunc("/articles",getAll).Methods("GET")

	http.ListenAndServe(":8000", router)
}