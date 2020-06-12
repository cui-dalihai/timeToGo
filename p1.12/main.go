package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
)

var mu sync.Mutex
var count int
var palt = []color.Color{color.White, color.Black}
const (
	whid = 0
	bkid = 1
)


func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/count", counter)
	http.HandleFunc("/report", report)
	http.HandleFunc("/lsjo", lsjo)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func report(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s	%s	%s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {

	rps := r.URL.Query()
	//if !err {
	//	fmt.Fprintf(w, "parameters not found.\n")
	//	return
	//}
	mu.Lock()
	count ++
	mu.Unlock()
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)

	for k, v := range rps {

		fmt.Fprintf(w, "params[%s] = %q\n", k, v)
	}

}

func counter(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	fmt.Fprintf(w, "Count %d\n", count)
	mu.Unlock()
}

//func lsjo(out io.Writer) {
func lsjo(w http.ResponseWriter, r *http.Request) {

	qps := r.URL.Query()
	strcycl := qps["cycl"][0]

	icycl, err := strconv.Atoi(strcycl)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	cycl := float64(icycl)

	const (
		res = 0.001
		size = 300
		nfra = 64
		dely = 8
	)
	freq := rand.Float64() * 3.0
	anim := gif.GIF{LoopCount:nfra}
	phas := 0.0

	for i:= 0; i < nfra; i ++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		imag := image.NewPaletted(rect, palt)

		for t := 0.0; t < cycl*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phas)
			imag.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), bkid)
		}

		phas += 0.1
		anim.Delay = append(anim.Delay, dely)
		anim.Image = append(anim.Image, imag)
	}
	gif.EncodeAll(w, &anim)
}


