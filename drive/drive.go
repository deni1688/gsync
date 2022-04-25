package drive

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func getClient(config *oauth2.Config) *http.Client {
	tokenFile := os.Getenv("HOME") + "/.gsync/token.json"
	token, err := tokenFromFile(tokenFile)
	if err != nil {
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

	ch := make(chan string)
	srv := &http.Server{Addr: ":9999"}

	go startTempAuthServer(ch, srv)

	log.Println("Waiting for auth code...")
	code := <-ch
	log.Println("Received auth code!")

	stopTemporaryServer(srv)

	if err := cmd.Process.Kill(); err != nil {
		log.Printf("Unable to kill browser: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}

	return tok
}

func stopTemporaryServer(srv *http.Server) {
	if err := srv.Shutdown(context.TODO()); err != nil {
		log.Fatalf("Unable to shutdown server: %v", err)
	}
	log.Println("Auth server stopped")
}

func startTempAuthServer(ch chan string, srv *http.Server) {
	http.HandleFunc("/oauth2callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		fmt.Fprint(w, `
			<html>
				<head><title>Gsync OAuth2</title></head>
				<body>
					<p style="font: 15px Arial, sans-serif;">Ahoy! You have been authenticated! Closing this window...</p>
				</body>
			</html>`)
		ch <- code
	})

	_ = srv.ListenAndServe()
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving token.json file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func New() *drive.Service {
	b, err := ioutil.ReadFile(os.Getenv("HOME") + "/.gsync/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	service, err := drive.NewService(context.TODO(), option.WithHTTPClient(getClient(config)))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	return service
}
