package main

import (
	"fmt"
	"log"
	"mime"
	"net/http"

	"github.com/husobee/vestigo"
	"github.com/mpcjanssen/adocsrv/pkg/adoc"
)

//go:generate go-bindata -prefix "assets" -pkg main -o bindata.go assets/...

const css = "asciidoctor-default.css"

func main() {
	router := vestigo.NewRouter()

	router.Get("/favicon.ico", favIcon)
	router.Get("/browse/*", adoc.BrowseHandler)
	router.Get("/view/*", adoc.ViewHandler)
	router.Get("/reveal/*", adoc.RevealHandler)
	router.Get("/edit/*", adoc.EditHandler)
	router.Get("/assets/*", assetsHandler)

	addr := ":12345"
	log.Println("Listening for requests on ", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	assetName := vestigo.Param(r, "_name")
	asset, err := Asset(assetName)
	if err != nil {
		fmt.Fprintln(w, err)
	} else {
		w.Header().Add("Content-Type", mime.TypeByExtension(assetName))
		w.Write(asset)
	}
}

func favIcon(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "image/png")
	favicon, _ := Asset("favicon.png")
	w.Write(favicon)
}

func adocCss(w http.ResponseWriter, _ *http.Request) {
	bin, _ := Asset(css)
	w.Header().Add("Content-Type", "text/css")
	w.Write(bin)
}
