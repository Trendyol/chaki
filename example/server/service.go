package main

import (
	"fmt"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GenerateText(text string, rt int) (string, error) {
	str := ""
	for i := 0; i < rt; i++ {
		temp := text + "\n"
		str += temp
	}

	return str, nil
}

func (s *Service) GenerateTextWithMetadata(text string, rt int, metadata map[string]string) (string, error) {
	str := ""

	// Add headers and cookies information to the response
	str += "Metadata received:\n"
	for key, value := range metadata {
		str += fmt.Sprintf("%s: %s\n", key, value)
	}
	str += "\n"

	// Add repeated text
	for i := 0; i < rt; i++ {
		temp := text + "\n"
		str += temp
	}

	return str, nil
}
