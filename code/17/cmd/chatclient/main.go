package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type loginResponse struct {
	Token string `json:"token"`
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		fmt.Println("\nexit")
		os.Exit(0)
	}()

	baseURL := "http://localhost:8082"
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("username: ")
	user, _ := reader.ReadString('\n')
	fmt.Print("password: ")
	pass, _ := reader.ReadString('\n')
	user = strings.TrimSpace(user)
	pass = strings.TrimSpace(pass)

	token, err := login(ctx, baseURL, user, pass)
	if err != nil {
		fmt.Println("login failed:", err)
		return
	}
	fmt.Println("login success.")

	for {
		if ctx.Err() != nil {
			return
		}
		fmt.Print("you> ")
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("read error:", err)
			return
		}
		msg = strings.TrimSpace(msg)
		if msg == "" {
			continue
		}

		fmt.Println("---- sending ----")
		done := make(chan struct{})
		go func() {
			defer close(done)
			if err := streamChat(ctx, baseURL, token, msg); err != nil {
				fmt.Println("\nchat error:", err)
			}
		}()

		select {
		case <-ctx.Done():
			fmt.Println("\nexit")
			return
		case <-done:
		}
	}
}

func login(ctx context.Context, baseURL, user, pass string) (string, error) {
	body, _ := json.Marshal(map[string]string{
		"username": user,
		"password": pass,
	})
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("status %d: %s", resp.StatusCode, string(b))
	}
	var lr loginResponse
	if err := json.NewDecoder(resp.Body).Decode(&lr); err != nil {
		return "", err
	}
	if lr.Token == "" {
		return "", fmt.Errorf("empty token")
	}
	return lr.Token, nil
}

func streamChat(ctx context.Context, baseURL, token, msg string) error {
	payload, _ := json.Marshal(map[string]string{
		"message": msg,
	})
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/chat", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(b))
	}

	reader := bufio.NewReader(resp.Body)
	fmt.Println("model> (streaming)")
	var printed bool
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if printed {
					fmt.Println()
				}
				return nil
			}
			return err
		}
		if strings.HasPrefix(line, "data: ") {
			chunk := strings.TrimPrefix(line, "data: ")
			fmt.Print(chunk)
			printed = true
		}
	}
}
