package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	Country string `json:"country,omitempty"`
}

type User struct {
	Name     string  `json:"name"`
	Password string  `json:"password"`
	Email    string  `json:"email"`
	Address  Address `json:"address"`
}

var users = []User{
	{
		Name:     "Bob",
		Password: "secret",
		Email:    "bob@golang.org",
		Address: Address{
			City:    "Lyon",
			Street:  "15 rue hade",
			Country: "France",
		},
	},
	{
		Name:     "Ginette",
		Password: "topsecret",
		Email:    "ginette@golang.org",
		Address: Address{
			City:   "Brest",
			Street: "12 rue du lavoir",
		},
	},
}

type PasswordJsonBody struct {
	UserIndex         int    `json:"user_index"`
	OldPassword       string `json:"old_password"`
	NewPassword       string `json:"new_password"`
	NewPasswordRepeat string `json:"new_password_repeat"`
}

func user(w http.ResponseWriter, r *http.Request) {
	//	encodage du tableau User en Json
	b, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func updatePassword(w http.ResponseWriter, r *http.Request) {
	var p PasswordJsonBody
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Parsed struct %v\n", p)

	if p.UserIndex < 0 || p.UserIndex > len(users)-1 {
		msg := fmt.Sprintf("Invalide index. got user_index=%v. valid range=[0,%v]", p.UserIndex, len(users)-1)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	u := users[p.UserIndex]

	if u.Password != p.OldPassword {
		http.Error(w, "Old password do not match", http.StatusBadRequest)
		return
	}

	if p.NewPassword != p.NewPasswordRepeat {
		http.Error(w, "New password do not match", http.StatusBadRequest)
		return
	}

	u.Password = p.NewPassword
	fmt.Fprintf(w, "Password updated")
}

func main() {
	http.HandleFunc("/users", user)
	http.HandleFunc("/update_password", updatePassword)

	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		log.Fatal(err)
	}
}
