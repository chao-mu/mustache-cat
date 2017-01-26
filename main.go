package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/unrolled/render"
)

type CatHandler struct {
	pngs [][]byte
	rend *render.Render
}

func NewCatHandler(imageDir string) (*CatHandler, error) {
	c := &CatHandler{rend: render.New()}
	files, err := ioutil.ReadDir(imageDir)
	if err != nil {
		return nil, err
	}

	for _, fileInfo := range files {
		path := fmt.Sprintf("%s/%s", imageDir, fileInfo.Name())
		if strings.HasSuffix(path, ".png") {
			err := c.addPNG(path)
			if err != nil {
				return nil, err
			}
		}
	}

	return c, nil
}

func (c *CatHandler) Handler() http.Handler {
	m := http.NewServeMux()

	m.HandleFunc("/", c.serveIndex)
	m.HandleFunc("/image.png", c.serveImage)

	return m
}

func (c *CatHandler) addPNG(path string) error {
	png, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	c.pngs = append(c.pngs, png)

	return nil
}

func (c *CatHandler) randCat() ([]byte, error) {
	if len(c.pngs) == 0 {
		return nil, errors.New("no cat PNGs loaded")
	}

	return c.pngs[rand.Int()%len(c.pngs)], nil
}

func (c *CatHandler) serveIndex(w http.ResponseWriter, req *http.Request) {
	page := map[string]interface{}{
		"ImageID": rand.Int(),
	}

	c.rend.HTML(w, http.StatusOK, "index", page)
}

func (c *CatHandler) serveImage(w http.ResponseWriter, req *http.Request) {
	png, err := c.randCat()
	if err != nil {
		Errorf(w, "Unable to retrieve random cat: %v", err)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(png)
}

func Errorf(w http.ResponseWriter, pat string, args ...interface{}) {
	http.Error(w, fmt.Sprintf(pat, args...), http.StatusInternalServerError)
}

func main() {
	h, err := NewCatHandler("images")
	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe("0.0.0.0:80", h.Handler()); err != nil {
		log.Fatal(err)
	}
}
