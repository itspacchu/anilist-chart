package api

import (
	"image"
	"image/jpeg"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gorilla/mux"
	"github.com/itspacchu/anilist-chart/processing"
)

func Run() error {
	r := mux.NewRouter()
	log.Printf("Running server on 0.0.0.0:3000")
	r.HandleFunc("/user/{username}/{type}", handleAPI)
	if err := http.ListenAndServe("0.0.0.0:3000", r); err != nil {
		return err
	}
	return nil
}

func handleAPI(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.Method, r.URL, r.RemoteAddr)
	vars := mux.Vars(r)
	activity_type := ""
	if mintype, ok := vars["type"]; ok {
		if mintype == "manga" {
			activity_type = "MANGA_LIST"
		} else {
			activity_type = "ANIME_LIST"
		}
	}
	if username, ok := vars["username"]; ok {
		var img *image.RGBA = processing.ProcessChart(username, 7, activity_type)
		w.Header().Set("Content-Type", "image/jpeg")
		if err := jpeg.Encode(w, img, &jpeg.Options{Quality: 90}); err != nil {
			http.Error(w, "Failed to encode JPEG", http.StatusInternalServerError)
		}
	}

}
