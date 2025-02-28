package ui

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

type resources struct {
	font text.Face
}

func LoadResources() (*resources, error) {
	fnt, err := loadFont(16)
	if err != nil {
		return nil, err
	}
	return &resources{
		font: fnt,
	}, nil
}

func loadFont(size float64) (text.Face, error) {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}
