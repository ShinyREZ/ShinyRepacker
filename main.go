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

const prefix = "unpacked"

func main() {
	var mode, file string

	flag.StringVar(&mode, "mode", "repack", "Work mode of ShinyRepacker")
	flag.StringVar(&file, "file", "parts_text", "File name without extension")
	flag.Parse()

	desc := LoadDescribeFile(file)

	switch mode {
	case "unpack":
		img := LoadImage(desc.Meta.Image)

		for name, frame := range desc.Frames {
			name = prefix + "/" + name
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
			SaveFile(name, result)
		}
	case "repack":
		result := image.NewNRGBA(image.Rect(0, 0, desc.Meta.Size.Width, desc.Meta.Size.Height))

		for name, frame := range desc.Frames {
			name = prefix + "/" + name
			img := LoadImage(name)

			f := frame.Frame

			// rotate: true means the width and height are exchanged

			if frame.Trimmed {
				// Trim before rotate
				triRect := image.Rect(0, 0, f.Width, f.Height)
				tri := image.NewNRGBA(triRect)
				draw.Draw(tri, triRect, img, image.Point{X: frame.SpriteSourceSize.X, Y: frame.SpriteSourceSize.Y}, draw.Over)
				img = tri
			}

			if frame.Rotated {
				rot := image.NewNRGBA(image.Rect(0, 0, img.Bounds().Dy(), img.Bounds().Dx()))
				for x := img.Bounds().Min.Y; x < img.Bounds().Max.Y; x++ {
					for y := img.Bounds().Max.X - 1; y >= img.Bounds().Min.X; y-- {
						rot.Set(img.Bounds().Max.Y-x, y, img.At(y, x))
					}
				}
				img = rot
			}

			if frame.Rotated {
				f.Width, f.Height = f.Height, f.Width
			}

			draw.Draw(
				result,
				image.Rect(f.X, f.Y, f.X+f.Width, f.Y+f.Height),
				img,
				image.Point{},
				draw.Over,
			)
		}

		SaveFile(file+".repack.png", result)
	}
}

func LoadImage(name string) image.Image {
	f, err := os.Open(name)
	if err != nil {
		log.Panic(err)
	}

	img, err := png.Decode(f)
	if err != nil {
		log.Panic(err)
	}
	return img
}

func SaveFile(name string, img image.Image) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		log.Panic(err)
	}

	folders := strings.Split(name, "/")
	if len(folders) > 1 {
		err = os.MkdirAll(strings.Join(folders[:len(folders)-1], "/"), 0777)
		if err != nil {
			log.Panic(err)
		}
	}

	err = ioutil.WriteFile(name, buf.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}
