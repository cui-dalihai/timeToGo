package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"log"
	"math"
	"math/cmplx"
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
	http.HandleFunc("/sufce", sufce)
	http.HandleFunc("/mdbrt", mdbrt)
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

const (
	width, height = 600, 320
	cells         = 100
	xyrange       = 30.0
	xyscale       = width / 2 / xyrange
	zscale        = height * 0.4
	angle         = math.Pi / 6
)
var sin30, cos30 = math.Sin(angle), math.Cos(angle)

func sufce(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "image/svg+xml")

	fmt.Fprintf(w, "<svg xmlns='http://www.w3.org/2000/svg' " +
		"style='stroke: grey; fill: white; stroke-width: 0.7' " +
		"width='%d' height='%d'>", width, height)

	for i := 0; i < cells; i ++ {
		for j := 0; j < cells; j ++ {
			ax, ay := corner(i+1, j)
			bx, by := corner(i, j)
			cx, cy := corner(i, j+1)
			dx, dy := corner(i+1, j+1)
			fmt.Fprintf(w, "<polygon points='%g,%g %g,%g %g,%g %g,%g'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	fmt.Fprintf(w, "</svg>")
}

func corner(i, j int) (float64, float64) {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	z := f(x, y)

	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y)
	return math.Sin(r) / r
}


func mdbrt(w http.ResponseWriter, r *http.Request) {
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height          = 1024, 1024
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py ++ {
		y := float64(py) / height * (ymax-ymin) + ymin
		for px := 0; px < width; px ++ {
			x := float64(px) / width * (xmax-xmin) + xmin
			z := complex(x, y)
			img.Set(px, py, mandelbrot(z))
		}
	}
	png.Encode(w, img)
}

func mandelbrot(z complex128) color.Color {
	const (
		iters = 200
		contr = 15
	)
	var v complex128
	for n := uint8(0); n < iters; n ++ {
		v = v * v + z
		if cmplx.Abs(v) > 2 {
			return color.Gray{ 255 - contr*n }
		}
	}
	return color.Black
}



