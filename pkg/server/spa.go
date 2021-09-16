package server

import (
	"embed"
	"github.com/gorilla/mux"
	"github.com/samvaughton/wpcommand/v2/pkg/config"
	log "github.com/sirupsen/logrus"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

type SpaHandler struct {
	IndexPath  string
	StaticPath string
	Embedded   *embed.FS
}

func (h SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	urlAbs, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path := filepath.Join(h.StaticPath, urlAbs)

	if h.Embedded == nil {

		// live spa

		// check whether a file exists at the given path
		_, err = os.Stat(path)
		if os.IsNotExist(err) {
			// file does not exist, serve index.html
			http.ServeFile(w, r, filepath.Join(h.StaticPath, h.IndexPath))
			return
		} else if err != nil {
			// if we got an error (that wasn't that the file doesn't exist) stating the
			// file, return a 500 internal server error and stop
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// otherwise, use http.FileServer to serve the static dir
		http.FileServer(http.Dir(h.StaticPath)).ServeHTTP(w, r)

	} else {

		// embedded spa
		fsys, err := fs.Sub(h.Embedded, h.StaticPath)
		if err != nil {
			panic(err)
		}

		var staticFS = http.FS(fsys)
		httpFs := http.FileServer(staticFS)

		file, err := staticFS.Open(urlAbs)

		if err != nil {
			log.Debug(err)
		}

		// no file found serve from index.html
		if file == nil {
			data, err := fs.ReadFile(fsys, h.IndexPath)

			_, err = w.Write(data)

			if err != nil {
				panic("could not serve index file via embedded FS")
			}

			return
		}

		// serve static files now
		httpFs.ServeHTTP(w, r)
	}

}

func SetupSpa(router *mux.Router, staticFiles *embed.FS) {
	log.Info("serving ui at /")

	spa := SpaHandler{StaticPath: "app/public", IndexPath: "index.html"}

	if config.Config.EmbedStaticFiles {
		spa = SpaHandler{StaticPath: "app/public", Embedded: staticFiles, IndexPath: "index.html"}
	}

	router.PathPrefix("/").Handler(spa)
}
