package admin

import (
	"context"
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
	"nut/gen/proto"
	"nut/internal"
)

var ENV = "debug"

//go:embed build
var embeddedFiles embed.FS

func ServeEmbeddedUI() error {
	// Get the build subdirectory as the
	// root directory so that it can be passed
	// to the http.FileServer
	fsys, err := fs.Sub(embeddedFiles, "build")
	if err != nil {
		return err
	}

	http.Handle("/", http.FileServer(http.FS(fsys)))
	return nil
}

func AddRoutes(dao internal.Persistance) {
	http.HandleFunc("/admin/tasks.list", listTasks(dao))
	http.HandleFunc("/admin/login", login)
}

func login(w http.ResponseWriter, r *http.Request) {
	// DANGER: get users from db and authenticate properly
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	json.NewEncoder(w).Encode(proto.DoneReply{Ok: true, Message: "Successfully logged in"})
}

// APIs

func listTasks(dao internal.Persistance) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks, err := dao.GetTasks(context.Background())
		if err != nil {
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(proto.DoneReply{Ok: false, Message: "Error while listing tasks"})
			return
		}

		json.NewEncoder(w).Encode(tasks)
	}
}
