package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type User struct {
	ID          string `json:"id"`
	Credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"credentials"`
}

func credentialCheck() {

	// Mengambil data dari API
	response, err := http.Get("http://localhost:8080/")

	if err != nil {
		fmt.Println("Gagal mengakses API:", err)
		return
	}

	defer response.Body.Close()

	// Menyimpan data JSON
	var users []User

	err = json.NewDecoder(response.Body).Decode(&users)

	if err != nil {
		fmt.Println("Gagal membaca data API:", err)
		return
	}

	// Input username
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Username: ")
	scanner.Scan()
	username := strings.TrimSpace(scanner.Text())

	// Input password
	fmt.Print("Password: ")
	scanner.Scan()
	password := strings.TrimSpace(scanner.Text())

	// Proses login
	loginBerhasil := false

	for _, user := range users {

		if user.Credentials.Username == username &&
			user.Credentials.Password == password {

			loginBerhasil = true

			fmt.Println("Login berhasil!")
			fmt.Println("User ID:", user.ID)

			break
		}
	}

	if !loginBerhasil {
		fmt.Println("Username atau password salah.")
	}
}
