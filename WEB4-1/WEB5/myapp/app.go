package myapp

import (
	"fmt"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"encoding/json"
	"strconv"
)

type User struct{
	ID 		  	int `json:"id"`
	FirstName 	string `json:"first_name"`
	LastName 	string `json:"last_name"`
	Email 		string  `json:"email"`
	CreatedAt 	time.Time `json: "created_at"`
}


type UpdateUser struct{
	ID 		  	int `json:"id"`
	FirstName 	string `json:"first_name"`
	LastName 	string `json:"last_name"`
	Email 		string  `json:"email"`
	CreatedAt 	time.Time `json: "created_at"`
}

var userMap map[int]*User
var lastID int

func indexHandler(w http.ResponseWriter, r *http.Request) {
	
	fmt.Fprint(w, "Hello World")
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	if len(userMap) == 0 {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No Users")
		return
	}

	users := []*User{}
	for _, u := range userMap {
		users = append(users, u)
	}

	data, _ := json.Marshal(users)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(data))
}

func getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return 
	}

	user, ok := userMap[id]
	if !ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User Id:", id)
		return
	}

	// user := new(User)
	// user.ID = 2
	// user.FirstName = "tucker"

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))
	// vars := mux.Vars(r) //request넣어주면 알아서 파싱해중
	// fmt.Fprint(w, "User Id:", vars["id"])
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err :=	json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	//Created User
	lastID++
	user.ID = lastID
	user.CreatedAt = time.Now()
	userMap[user.ID] = user

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))
}

func deleteUserInfoHandler(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"]) //integer가 아닐수도 있으니까 atoi 함 사용자가 보낸게 int가 아닐수도 ..

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return 
	}

	_, ok := userMap[id]
	if !ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User Id:", id)
		return
	}

	delete(userMap, id)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Deleted User Id:", id)

}

func updateUserHandler(w http.ResponseWriter, r *http.Request){
	updateUser := new(User)
	err :=	json.NewDecoder(r.Body).Decode(updateUser)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	user, ok := userMap[updateUser.ID]
	if !ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User Id:", updateUser.ID)
		return
	}

	if updateUser.FirstName != "" {
		user.FirstName = updateUser.FirstName
	}


	if updateUser.LastName != "" {
		user.LastName = updateUser.LastName
	}

	if updateUser.Email != "" {
		user.Email = updateUser.Email
	}

	//userMap[updateUser.ID] = updateUser
	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))
}

//make a new myapp handler
func NewHandler() http.Handler {

	userMap = make(map[int]*User)
	lastID = 0
	mux := mux.NewRouter()
	//mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/users", usersHandler).Methods("GET")
	mux.HandleFunc("/users", createUserHandler).Methods("POST")
	mux.HandleFunc("/users", updateUserHandler).Methods("PUT")

	mux.HandleFunc("/users/{id:[0-9]+}", getUserInfoHandler).Methods("GET")
	mux.HandleFunc("/users/{id:[0-9]+}", deleteUserInfoHandler).Methods("DELETE")

	return mux
}