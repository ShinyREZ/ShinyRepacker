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
	"path"
	"strings"
)

func main() {
	var mode, file, imgInput, prefix string

	flag.StringVar(&mode, "mode", "unpack", "Work mode, unpack or repack")
	flag.StringVar(&file, "file", "", "Json file name")
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
		imgInput = path.Dir(file) + "/" + imgInput

		img := LoadImage(imgInput)

		for name, frame := range desc.Frames {
			if prefix != "" {
				name = path.Dir(file) + "/" + prefix + "/" + name
			}
			f := frame.Frame

			// rotate: true means the width and height are exchanged
			if frame.Rotated {
				f.Width, f.Height = f.Height, f.Width
			}

			recRect := image.Rect(0, 0, f.Width+2, f.Height+2)
			rec := image.NewNRGBA(recRect)
			draw.Draw(rec, recRect, img, image.Point{X: f.X - 1, Y: f.Y - 1}, draw.Over)

			if frame.Rotated {
				f.Width, f.Height = f.Height, f.Width
				rot := image.NewNRGBA(image.Rect(0, 0, f.Width+2, f.Height+2))
				for x := 0; x < f.Width+2; x++ {
					for y := 0; y < f.Height+2; y++ {
						rot.Set(x, y, rec.At(f.Height+2-1-y, x)) // Important: Height -1 here
					}
				}
				rec = rot
			}

			var result image.Image = rec
			if frame.Trimmed {
				dst := image.NewNRGBA(image.Rect(0, 0, frame.SourceSize.Width+2, frame.SourceSize.Height+2))
				draw.Draw(
					dst,
					image.Rect(
						frame.SpriteSourceSize.X-1,
						frame.SpriteSourceSize.Y-1,
						frame.SpriteSourceSize.X+frame.SpriteSourceSize.Width+1,
						frame.SpriteSourceSize.Y+frame.SpriteSourceSize.Height+1,
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
			fr := frame.SpriteSourceSize

			if frame.Trimmed {
				// Trim before rotate
				triRect := image.Rect(0, 0, f.Width+2, f.Height+2)
				tri := image.NewNRGBA(triRect)
				draw.Draw(tri, triRect, img, image.Point{X: fr.X - 1, Y: fr.Y - 1}, draw.Over)
				img = tri
			}

			if frame.Rotated {
				f.Width, f.Height = f.Height, f.Width
				rot := image.NewNRGBA(image.Rect(0, 0, f.Width+2, f.Height+2))
				for x := 0; x < f.Width+2; x++ {
					for y := 0; y < f.Height+2; y++ {
						rot.Set(x, y, img.At(y, f.Width+2-1-x)) // Important: Width -1 here
					}
				}
				img = rot
			}

			draw.Draw(
				result,
				image.Rect(f.X-1, f.Y-1, f.X+f.Width+1, f.Y+f.Height+1),
				img,
				image.Point{},
				draw.Src,
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
