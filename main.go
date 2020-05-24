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
	var mode, file, imgInput, prefix string

	flag.StringVar(&mode, "mode", "unpack", "Work mode, unpack or repack")
	flag.StringVar(&file, "file", "", "Json file name without extension")
	flag.StringVar(&imgInput, "image", "", "Image file name, used as input in unpack mode or as output in repack mode")
	flag.StringVar(&prefix, "prefix", "unpacked", "Prefix path of exported files")
	flag.Parse()

	if file == "" {
		log.Fatal("No file input provided")
		return
	}

	desc := LoadDescribeFile(file)

	switch mode {
	case "unpack":
		if imgInput == "" {
			imgInput = desc.Meta.Image
		}

		img := LoadImage(imgInput)

		for name, frame := range desc.Frames {
			if prefix != "" {
				name = prefix + "/" + name
			}
			f := frame.Frame

			// rotate: true means the width and height are exchanged
			if frame.Rotated {
				f.Width, f.Height = f.Height, f.Width
			}

			recRect := image.Rect(0, 0, f.Width, f.Height)
			rec := image.NewNRGBA(recRect)
			draw.Draw(rec, recRect, img, image.Point{X: f.X, Y: f.Y}, draw.Over)

			if frame.Rotated {
				f.Width, f.Height = f.Height, f.Width
				rot := image.NewNRGBA(image.Rect(0, 0, f.Width, f.Height))
				for x := 0; x < f.Width; x++ {
					for y := 0; y < f.Height; y++ {
						rot.Set(x, y, rec.At(f.Height-1-y, x)) // Important: Height -1 here
					}
				}
				rec = rot
			}

			var result image.Image = rec
			if frame.Trimmed {
				dst := image.NewNRGBA(image.Rect(0, 0, frame.SourceSize.Width, frame.SourceSize.Height))
				draw.Draw(
					dst,
					image.Rect(
						frame.SpriteSourceSize.X,
						frame.SpriteSourceSize.Y,
						frame.SourceSize.Width,
						frame.SourceSize.Height,
					),
					rec, image.Point{}, draw.Over,
				)
				result = dst
			}

			// Save file
			SaveFile(name, result)
		}
	case "repack":
		result := image.NewNRGBA(image.Rect(0, 0, desc.Meta.Size.Width, desc.Meta.Size.Height))

		for name, frame := range desc.Frames {
			if prefix != "" {
				name = prefix + "/" + name
			}

			img := LoadImage(name)
			f := frame.Frame

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

		// Save file
		var fileName string
		if imgInput == "" {
			fileName = file + ".repack.png"
		} else {
			fileName = imgInput
		}
		SaveFile(fileName, result)
	default:
		flag.Usage()
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
