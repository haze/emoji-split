package main

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"os"
)

func c(e error) {
	if e != nil {
		panic(e)
	}
}

const (
	pixelChunkSize int = 16
)

func main() {
	img, err := os.Open("test.gif")
	c(err)
	dat, err := gif.DecodeAll(img)
	c(err)
	err = os.Mkdir("test", os.ModePerm)
	c(err)
	bounds := dat.Image[0].Rect
	frameMap := make(map[image.Point][]*image.Paletted)
	for y := 0; y <= bounds.Max.Y; y += pixelChunkSize {
		for x := 0; x <= bounds.Max.X; x += pixelChunkSize {
			fmt.Printf("x{%d}, y{%d}\n", x, y)
			pt := image.Pt(x, y)
			for frameIndex := range dat.Image {
				frame := dat.Image[frameIndex]
				rect := image.Rect(x, y, x+pixelChunkSize, y+pixelChunkSize)
				pal := image.NewPaletted(rect, palette.WebSafe)
				z := frame.SubImage(image.Rect(x, y, x+pixelChunkSize, y+pixelChunkSize))
				draw.Draw(pal, rect, z, pal.Bounds().Min, draw.Over)
				frameMap[pt] = append(frameMap[pt], pal)
			}
		}
	}
	count := 0
	fmt.Println(len(frameMap))
	for _, v := range frameMap {
		file, err := os.Create(fmt.Sprintf("test/chunk-%d.gif", count))
		c(err)
		g := gif.GIF{
			Delay:           dat.Delay,
			LoopCount:       dat.LoopCount,
			Disposal:        dat.Disposal,
			Config:          dat.Config,
			BackgroundIndex: dat.BackgroundIndex,
			Image:           v,
		}
		gif.EncodeAll(file, &g)
		count++
	}
}
