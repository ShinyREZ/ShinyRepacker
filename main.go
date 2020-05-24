package main

import (
	"bytes"
	"flag"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	var mode, file string

	flag.StringVar(&mode, "mode", "unpack", "Work mode of ShinyRepacker")
	flag.StringVar(&file, "file", "parts_text", "File name without extension")
	flag.Parse()

	desc := LoadDescribeFile(file)

	switch mode {
	case "unpack":
		f, err := os.Open(desc.Meta.Image)
		if err != nil {
			log.Panic(err)
		}

		img, err := png.Decode(f)
		if err != nil {
			log.Panic(err)
		}

		for name, frame := range desc.Frames {
			name = strings.ReplaceAll(name, "/", "_")

			f := frame.Frame

			// rotate: true means the width and height are exchanged
			if frame.Rotated {
				f.Width, f.Height = f.Height, f.Width
			}

			recRect := image.Rect(0, 0, f.Width, f.Height)

			rec := image.NewNRGBA(recRect)
			draw.Draw(rec, recRect, img, image.Point{X: f.X, Y: f.Y}, draw.Over)

			if frame.Rotated {
				rot := image.NewNRGBA(image.Rect(0, 0, f.Height, f.Width))
				for x := 0; x < f.Height; x++ {
					for y := f.Width - 1; y >= 0; y-- {
						rot.Set(x, f.Width-y, rec.At(y, x))
					}
				}
				rec = rot
			}

			var result image.Image
			if frame.Trimmed {
				dst := image.NewNRGBA(image.Rect(0, 0, frame.SourceSize.Width, frame.SourceSize.Height))
				draw.Draw(
					dst,
					image.Rect(
						frame.SpriteSourceSize.X,
						frame.SpriteSourceSize.Y,
						frame.SpriteSourceSize.X+frame.SpriteSourceSize.Width,
						frame.SpriteSourceSize.Y+frame.SpriteSourceSize.Height,
					),
					rec, image.Point{}, draw.Over,
				)
				result = dst
			} else {
				result = rec
			}

			// Save file
			var buf bytes.Buffer
			err := png.Encode(&buf, result)
			if err != nil {
				log.Panic(err)
			}

			err = ioutil.WriteFile(name, buf.Bytes(), 0644)
			if err != nil {
				log.Panic(err)
			}
		}
	case "repack":
	}
}
