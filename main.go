package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w,r)
	})
}



func main () {
var apiCfg apiConfig
mux := http.NewServeMux()
fileHandler := http.FileServer(http.Dir("./"))
mux.Handle("/app/",apiCfg.middlewareMetricsInc(http.StripPrefix("/app",fileHandler)))
mux.HandleFunc("GET /api/healthz",healthHandler)
mux.HandleFunc("GET /admin/metrics",apiCfg.hitsHandler)
mux.HandleFunc("POST /admin/reset",apiCfg.resethitsHandler)
mux.HandleFunc("POST /api/validate_chirp", validatechirpHandler)
srv := &http.Server{
	Addr: ":8081",
	Handler: mux,
}
err := srv.ListenAndServe()
if err != nil  {
	 fmt.Printf("error when starting server:%v ",err)
}
}