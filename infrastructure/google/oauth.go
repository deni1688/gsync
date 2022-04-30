package google

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func getClient(config *oauth2.Config) *http.Client {
	tokenFile := os.Getenv("HOME") + "/.gsync/token.json"
	token, err := tokenFromFile(tokenFile)

	if err != nil || !token.Valid() {
		token = getTokenFromWeb(config)
		saveToken(tokenFile, token)
	}

	return config.Client(context.Background(), token)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	cmd := exec.Command("xdg-open", authURL)
	if err := cmd.Start(); err != nil {
		log.Fatalf("Unable to open browser: %v", err)
	}

	codeCh := make(chan string)
	srv := &http.Server{Addr: ":9999"}

	go startTempAuthServer(codeCh, srv)

	log.Println("Waiting for auth code...")
	code := <-codeCh
	log.Println("Received auth code!")

	stopTempAuthServer(srv)

	if err := cmd.Process.Kill(); err != nil {
		log.Printf("Unable to kill browser: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}

	return tok
}

func stopTempAuthServer(srv *http.Server) {
	if err := srv.Shutdown(context.TODO()); err != nil {
		log.Fatalf("Unable to shutdown server: %v", err)
	}
	log.Println("Auth server stopped")
}

func startTempAuthServer(codeCh chan string, srv *http.Server) {
	http.HandleFunc("/oauth2callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		fmt.Fprint(w, `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Gsync Auth</title>
    <link rel="stylesheet" href="https://bootswatch.com/5/lux/bootstrap.min.css">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.8.1/font/bootstrap-icons.css">
    <style>
        html, body {
            height: 100%;
        }
    </style>
</head>
<body class="bg-dark d-flex flex-column justify-content-center align-items-center">
<div class="col-3">
    <div class="card card-body text-center">
        <i class="bi bi-check-circle-fill text-success" style="font-size: 50px"></i>
        <h4>Auth Success!</h4>
        <p>You can close this page now.</p>
    </div>
</div>
</body>
</html>`)
		codeCh <- code
	})

	_ = srv.ListenAndServe()
}

func tokenFromFile(path string) (*oauth2.Token, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	tok := new(oauth2.Token)
	err = json.NewDecoder(file).Decode(tok)

	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving token.json file to: %s\n", path)

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer file.Close()

	if err = json.NewEncoder(file).Encode(token); err != nil {
		log.Fatalf("Unable to encode oauth token: %v", err)
	}
}
