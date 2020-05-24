package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type DescribeFile struct {
	Frames map[string]DescribeFrame `json:"frames"`
	Meta   DescribeMeta             `json:"meta"`
}

type XYWH struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"w"`
	Height int `json:"h"`
}

type WH struct {
	Width  int `json:"w"`
	Height int `json:"h"`
}

type DescribeFrame struct {
	Frame            XYWH `json:"frame"`
	Rotated          bool `json:"rotated"`
	Trimmed          bool `json:"trimmed"`
	SpriteSourceSize XYWH `json:"spriteSourceSize"`
	SourceSize       WH   `json:"sourceSize"`
}

type DescribeMeta struct {
	App     string `json:"app"`
	Version string `json:"version"`
	Image   string `json:"image"`
	Format  string `json:"format"`
	Size    struct {
		Width  int `json:"w"`
		Height int `json:"h"`
	} `json:"size"`
	Scale       string `json:"scale"`
	SmartUpdate string `json:"smartupdate"`
}

func LoadDescribeFile(name string) DescribeFile {
	var result DescribeFile

	data, err := ioutil.ReadFile(name + ".json")
	if err != nil {
		log.Panic(err)
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Panic(err)
	}

	return result
}
