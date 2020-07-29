package main

import (
	"errors"
)

type Asset struct {
	ID     string
	Format string
	Path   string
}

func GetAsset(id string) (Asset, error) {
	return Asset{}, errors.New("asset does not exist")
}

