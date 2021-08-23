package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
)

func (app *App) resize(w http.ResponseWriter, r *http.Request) {
	var details struct {
		UserID UserID
		Width  uint
		Height uint
	}
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"details": "expect JSON"})
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&details); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"details": "invalid JSON"})
		log.Printf("decoding json error: %s", err)
		return
	}

	app.lock.Lock()
	infos, ok := app.users[details.UserID]
	app.lock.Unlock()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"details": "invalid user id"})
		log.Printf("invalid user id: %q", details.UserID)
		return
	}

	execID := infos.ExecID

	if err := app.client.ContainerExecResize(context.Background(), string(execID), types.ResizeOptions{
		Width:  details.Width,
		Height: details.Height,
	}); err != nil {
		panic(err)
	}
}