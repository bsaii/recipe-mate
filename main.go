package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Difficulty string

const (
	Easy     Difficulty = "Easy"
	Medium   Difficulty = "Medium"
	Hard     Difficulty = "Hard"
	Advanced Difficulty = "Advanced"
)

type Recipe struct {
	Id              string     `json:"id"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Ingredients     []string   `json:"ingredients"`
	Instructions    string     `json:"instructions"`
	CookingTime     int        `json:"cooking_time"`
	DifficultyLevel Difficulty `json:"difficulty_level"`
	Username        string     `json:"username"`
	CreatedAt       string     `json:"created_at"`
	UpdatedAt       string     `json:"updated_at"`
}

var recipes []Recipe

func getRecipes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

func delRecipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, recipe := range recipes {
		if recipe.Id == params["id"] {
			recipes = append(recipes[:index], recipes[index+1:]...)
			break
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(recipes)
}

func getRecipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, recipe := range recipes {
		if recipe.Id == params["id"] {
			json.NewEncoder(w).Encode(recipe)
			return
		}
	}
}

func addRecipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var recipe Recipe

	_ = json.NewDecoder(r.Body).Decode(&recipe)
	recipe.Id = strconv.Itoa(rand.Intn(1000000))
	recipes = append(recipes, recipe)
	json.NewEncoder(w).Encode(recipe)
}

func updateRecipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, recipe := range recipes {
		if recipe.Id == params["id"] {
			recipes = append(recipes[:index], recipes[index+1:]...)
			var recipe Recipe
			_ = json.NewDecoder(r.Body).Decode(&recipe)
			recipe.Id = params["id"]
			recipes = append(recipes, recipe)
			json.NewEncoder(w).Encode(recipe)
			return
		}
	}
}

func main() {
	r := mux.NewRouter()

	recipes = append(recipes, Recipe{Id: "1", Title: "Spaghetti Bolognese", Description: "Classic Italian pasta dish with meaty sauce", Ingredients: []string{"spaghetti", "ground beef", "tomato sauce", "onion", "garlic", "olive oil"}, Instructions: "1. Boil the spaghetti.\n2. Brown the beef with onions and garlic.\n3. Add tomato sauce and simmer.\n4. Serve sauce over cooked spaghetti.",
		CookingTime: 30, DifficultyLevel: Medium,
		Username: "chefmaster", CreatedAt: "2022-02-28", UpdatedAt: "2022-02-28"})

	recipes = append(recipes, Recipe{Id: "2",
		Title:           "Grilled Chicken Salad",
		Description:     "Healthy and delicious salad with grilled chicken",
		Ingredients:     []string{"chicken breast", "lettuce", "tomatoes", "cucumbers", "olives", "feta cheese"},
		Instructions:    "1. Grill the chicken breast.\n2. Chop vegetables and mix in a bowl.\n3. Slice grilled chicken and add to the salad.\n4. Drizzle with your favorite dressing.",
		CookingTime:     25,
		DifficultyLevel: Easy,
		Username:        "healthyeater", CreatedAt: "2022-03-05", UpdatedAt: "2022-03-05"})

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Recipe Mate!\n"))
	}).Methods("GET")
	r.HandleFunc("/recipes", getRecipes).Methods("GET")
	r.HandleFunc("/recipe", addRecipe).Methods("POST")
	r.HandleFunc("/recipe/{id}", getRecipe).Methods("GET")
	r.HandleFunc("/recipe/{id}", updateRecipe).Methods("PUT")
	r.HandleFunc("/recipe/{id}", delRecipe).Methods("DELETE")

	fmt.Printf("Starting server at port 8000\n")
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	log.Fatal(http.ListenAndServe(":8000", loggedRouter))
}
