package main

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
