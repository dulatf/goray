package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
)

// Eps : Minimal distance
const Eps = 1.0e-4
const maxSteps = 128
const numSamples = 16

// Light : Simple point light
type Light struct {
	position, color Vector
}

func clamp(s float64) uint8 {
	return uint8(math.Min(math.Max(255.0*s, 0.0), 255.0))
}

func save(filename string, w, h int, framebuffer []Vector) {
	f, _ := os.Create(filename)
	defer f.Close()
	wr := bufio.NewWriter(f)
	defer wr.Flush()
	fmt.Fprintf(wr, "P3\n%d %d\n255\n", w, h)
	for _, p := range framebuffer {
		r := clamp(p.x)
		g := clamp(p.y)
		b := clamp(p.z)
		fmt.Fprintf(wr, "%d %d %d\n", r, g, b)
	}
}

func opUnion(d1, d2 float64) float64 {
	return math.Min(d1, d2)
}
func opSub(d1, d2 float64) float64 {
	return math.Max(d1, -d2)
}

func sdSphere(p Vector, r float64) float64 {
	return p.Length() - r
}
func sdTorus(p Vector, r1, r2 float64) float64 {
	qx, qy := math.Sqrt(p.x*p.x+p.z*p.z)-r1, p.y
	return math.Sqrt(qx*qx+qy*qy) - r2
}

func signedDistance(p Vector) float64 {
	/*return opSub(opSub(opUnion(opUnion(opUnion(
	sdSphere(p, 1.5),
	sdSphere(Sub(p, Vector{2, 0, 0}), 1.0)),
	sdSphere(Sub(p, Vector{0, -100, 0}), 100.0)),
	sdTorus(Sub(p, Vector{0, 2, 0}), 1.0, 0.2)),
	sdTorus(Sub(p, Vector{0, 1, 0}), 1.0, 0.2)),
	sdSphere(Sub(p, Vector{-2, 0, -2}), 1.0))*/
	return opSub(sdSphere(p, 2.0),
		sdTorus(p, 2.0, 0.5))
}

func normal(pos Vector) Vector {
	delx := Vector{Eps, 0, 0}
	dely := Vector{0, Eps, 0}
	delz := Vector{0, 0, Eps}
	n := Vector{signedDistance(Add(pos, delx)) - signedDistance(Sub(pos, delx)),
		signedDistance(Add(pos, dely)) - signedDistance(Sub(pos, dely)),
		signedDistance(Add(pos, delz)) - signedDistance(Sub(pos, delz))}
	n.Normalize()
	return n
}
func lightContribution(pos Vector, light Light) Vector {
	n := normal(pos)
	v := Sub(light.position, pos)
	v.Normalize()
	hit, _ := sphereTrace(pos, v)
	if hit {
		return Vector{0.0, 0.0, 0.0}
	} else {
		return Mul(light.color, Dot(n, v))
	}
}
func shade(pos Vector) Vector {
	/*lights := [3]Light{Light{Vector{2.0, 5.0, -1.0}, Vector{0.3, 0.3, 0.3}},
	Light{Vector{-2.0, 5.0, -1.0}, Vector{0.3, 0.3, 0.3}},
	Light{Vector{0.0, 10.0, 0.0}, Vector{0.1, 0.15, 0.2}}}*/
	lights := [2]Light{Light{Vector{3.0, 2.0, -2.0}, Vector{1.0, 0.5, 0.4}},
		Light{Vector{-4.0, -1.0, -2.0}, Vector{0.2, 1.0, 0.5}}}
	col := Vector{0.0, 0.0, 0.0}
	for _, light := range lights {
		col = Add(col, lightContribution(pos, light))
	}
	return col
}

func sphereTrace(o, dir Vector) (bool, Vector) {
	t := Eps
	for i := 0; i < maxSteps; i++ {
		pos := Add(o, Mul(dir, t))
		d := signedDistance(pos)
		if math.Abs(d) < Eps*t {
			return true, pos
		}
		t = t + d
		if t > 100 {
			break
		}
	}
	return false, Vector{0.0, 0.0, 0.0}
}

func tracePixel(w, h int, x, y float64, fov float64) Vector {
	o := Vector{0, 0, -5.0}
	d := Vector{x + 0.5 - float64(w)/2.0,
		-(y + 0.5 - float64(h)/2.0),
		float64(h) / (2.0 * math.Tan(fov/2.0))}
	d.Normalize()
	hit, pos := sphereTrace(o, d)
	if hit {
		return shade(pos)
	} else {
		return Vector{0.34, 0.6, 0.8}
	}
}

func render(w, h int, framebuffer []Vector) {
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			col := Vector{0, 0, 0}
			for i := 0; i < numSamples; i++ {
				px := float64(x) - 0.5 + rand.Float64()
				py := float64(y) - 0.5 + rand.Float64()
				col = Add(col, tracePixel(w, h, px, py, math.Pi/3.0))
			}
			framebuffer[y*w+x] = Mul(col, 1.0/float64(numSamples))
		}
	}
}

func main() {
	fmt.Println("Hello, world y'all!")
	w, h := 640, 360
	framebuffer := make([]Vector, w*h)
	render(w, h, framebuffer)
	save("out.ppm", w, h, framebuffer)
}
