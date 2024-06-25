package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/go-chi/chi"
)

type User struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	Age     int     `json:"age"`
	Friends []int64 `json:"friends"`
}

type FileStore struct {
	FilePath string
	Users    map[int64]User
	NextID   int64
	mutex    sync.Mutex
}

var fileStore = FileStore{
	FilePath: "users.json",
	Users:    make(map[int64]User),
	NextID:   1,
}

func init() {
	if err := loadFromFile(&fileStore); err != nil {
		log.Fatalf("Failed to load data from file: %v", err)
	}
}

func main() {
	router := chi.NewRouter()
	router.Post("/create", createUser)
	router.Post("/make_friends", makeFriends)
	router.Delete("/user/{id}", deleteUser)
	router.Get("/friends/{id}", getFriends)
	router.Put("/user/{id}", updateAge)

	go func() {
		log.Println("Starting server 1 on port 8080...")
		log.Fatal(http.ListenAndServe(":8080", router))
	}()

	go func() {
		log.Println("Starting server 2 on port 8081...")
		log.Fatal(http.ListenAndServe(":8081", router))
	}()

	select {}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileStore.mutex.Lock()
	defer fileStore.mutex.Unlock()

	user.ID = fileStore.NextID
	fileStore.Users[user.ID] = user
	fileStore.NextID++

	if err := saveToFile(&fileStore); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"id": user.ID})
}

func makeFriends(w http.ResponseWriter, r *http.Request) {
	var data struct {
		SourceID int64 `json:"source_id"`
		TargetID int64 `json:"target_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileStore.mutex.Lock()
	defer fileStore.mutex.Unlock()

	sourceUser, ok := fileStore.Users[data.SourceID]
	if !ok {
		http.Error(w, "Source user not found", http.StatusNotFound)
		return
	}

	targetUser, ok := fileStore.Users[data.TargetID]
	if !ok {
		http.Error(w, "Target user not found", http.StatusNotFound)
		return
	}

	sourceUser.Friends = append(sourceUser.Friends, targetUser.ID)
	targetUser.Friends = append(targetUser.Friends, sourceUser.ID)

	if err := saveToFile(&fileStore); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s and %s are now friends\n", sourceUser.Name, targetUser.Name)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileStore.mutex.Lock()
	defer fileStore.mutex.Unlock()

	_, ok := fileStore.Users[userID]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	for _, user := range fileStore.Users {
		for i, id := range user.Friends {
			if id == userID {
				user.Friends = append(user.Friends[:i], user.Friends[i+1:]...)
				break
			}
		}
	}

	delete(fileStore.Users, userID)

	if err := saveToFile(&fileStore); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User with ID %d has been deleted\n", userID)
}

func getFriends(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileStore.mutex.Lock()
	defer fileStore.mutex.Unlock()

	user, ok := fileStore.Users[userID]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	friends := make([]string, len(user.Friends))
	for i, friendID := range user.Friends {
		friend, _ := fileStore.Users[friendID]
		friends[i] = friend.Name
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(friends)
}

func updateAge(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var data struct {
		NewAge int `json:"new_age"`
	}
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileStore.mutex.Lock()
	defer fileStore.mutex.Unlock()

	user, ok := fileStore.Users[userID]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	user.Age = data.NewAge
	fileStore.Users[userID] = user

	if err := saveToFile(&fileStore); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User's age has been updated to %d\n", user.Age)
}

func saveToFile(store *FileStore) error {
	file, err := os.Create(store.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(store); err != nil {
		return err
	}

	return nil
}

func loadFromFile(store *FileStore) error {
	file, err := os.Open(store.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(store); err != nil {
		return err
	}

	return nil
}
