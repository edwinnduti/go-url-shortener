/*
[*] Copyright © 2020
[*] Dev/Author ->  Edwin Nduti
[*] Description:
	A Golang and mongoDB Url-shortener similar to https://bit.ly
*/

package main

// libraries to use
import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
	//"github.com/urfave/negroni"
	"github.com/speps/go-hashids"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// handle our golang structure
type MyUrl struct {
	Id        primitive.ObjectID `bson:"_id",json:"id"`
	UrlID     string             `json:"urlid"`
	LongUrl   string             `json:"longurl"`
	ShortUrl  string             `json:"shorturl"`
	CreatedAt time.Time          `json:"createdat"`
	UpdatedAt time.Time          `json:"updatedat"`
}

// database and collection names are statically declared
const database, collection = "url-shortener", "urls"

// handle error
func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// HTTP /POST the / route will allow us to pass in a long URL and receive a short URL
func CreateEndpoint(w http.ResponseWriter, r *http.Request) {
	var url MyUrl
	err := json.NewDecoder(r.Body).Decode(&url)
	Check(err)

	client, err := CreateConnection()
	Check(err)

	c := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//create the shorturl
	url.Id = primitive.NewObjectID()
	hd := hashids.NewData()
	h, err := hashids.NewWithData(hd)
	Check(err)
	now := time.Now()
	url.UrlID, _ = h.Encode([]int{int(now.Unix())})
	url.ShortUrl = r.Host + "/" + url.UrlID
	url.CreatedAt = time.Now()
	_, err = c.InsertOne(ctx, url)
	Check(err)

	jsonData, err := json.Marshal(url)
	Check(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

// HTTP /GET /expand endpoint will allow us to pass in a short URL and receive a long URL
func ExpandEndpoint(w http.ResponseWriter, r *http.Request) {
	var url MyUrl
	err := json.NewDecoder(r.Body).Decode(&url)
	Check(err)
	shorturl := url.ShortUrl

	var result MyUrl

	client, err := CreateConnection()
	Check(err)

	c := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = c.FindOne(ctx, bson.M{"shorturl": shorturl}).Decode(&result)
	Check(err)

	jsonData, err := json.Marshal(result)
	Check(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

// HTTP /GET url data
func GetUrlDataHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	obj_id, err := primitive.ObjectIDFromHex(vars["id"])
	Check(err)

	var url MyUrl

	client, err := CreateConnection()
	Check(err)

	c := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = c.FindOne(ctx, bson.M{"_id": obj_id}).Decode(&url)
	Check(err)

	jsonData, err := json.Marshal(url)
	Check(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// HTTP /GET /{urlid} the /root endpoint will allow us to pass in a hash and be redirected to the long URL page
func RootEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	obj_id := vars["urlid"]

	var url MyUrl

	client, err := CreateConnection()
	Check(err)

	c := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = c.FindOne(ctx, bson.M{"urlid": obj_id}).Decode(&url)
	Check(err)

	// Redirect to long url
	http.Redirect(w, r, string(url.LongUrl), 301)
}

// HTTP /PUT user record /{id}
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	obj_id, err := primitive.ObjectIDFromHex(vars["id"])
	Check(err)

	url := MyUrl{}

	var updatedUrl MyUrl

	// Decode the incoming Data json
	err = json.NewDecoder(r.Body).Decode(&updatedUrl)
	Check(err)

	client, err := CreateConnection()
	Check(err)

	c := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = c.FindOne(ctx, bson.M{"_id": obj_id}).Decode(&url)
	Check(err)

	url.LongUrl = updatedUrl.LongUrl
	url.UpdatedAt = time.Now()

	// update longurl and show update time
	var update bson.M
	update = bson.M{
		"$set": bson.M{
			"longurl":   url.LongUrl,
			"updatedat": url.UpdatedAt,
		},
	}

	_, err = c.UpdateOne(ctx, bson.M{"_id": obj_id}, update)
	Check(err)

	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(url)
	Check(err)

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// HTTP /DELETE url record /{id}
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	obj_id, err := primitive.ObjectIDFromHex(vars["id"])
	Check(err)

	var url MyUrl

	client, err := CreateConnection()
	Check(err)

	c := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//err = c.FindOne(ctx,bson.M{"_id": obj_id}).Decode(&url)
	//Check(err)
	_, err = c.DeleteOne(ctx, bson.M{"_id": obj_id})
	Check(err)

	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(url)
	Check(err)

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func CreateConnection() (*mongo.Client, error) {
	// connect to mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	MongoURI := os.Getenv("MONGOURI")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		MongoURI,
	))
	Check(err)
	return client, nil
}

// Main function
func main() {
	/*
	   mgo.SetDebug(true)
	   mgo.SetLogger(log.New(os.Stdout,"err",6))

	   The above two lines are for debugging errors
	   that occur straight from accessing the mongo db
	*/

	//Register router{id}
	r := mux.NewRouter().StrictSlash(false)

	// API routes,handlers and methods
	r.HandleFunc("/", CreateEndpoint).Methods("POST")
	r.HandleFunc("/{id}", GetUrlDataHandler).Methods("GET")
	r.HandleFunc("/{id}", UpdateUserHandler).Methods("PUT")
	r.HandleFunc("/{id}", DeleteUserHandler).Methods("DELETE")
	r.HandleFunc("/expand", ExpandEndpoint).Methods("GET")
	r.HandleFunc("/{urlid}", RootEndpoint).Methods("GET")

	//Get port
	Port := os.Getenv("PORT")
	if Port == "" {
		Port = "8045"
	}

	// establish logger
	//n := negroni.Classic()
	//n.UseHandler(r)
	server := &http.Server{
		Handler: r, // n
		Addr:    ":" + Port,
	}
	log.Printf("Listening on PORT: %s", Port)
	server.ListenAndServe()
}
