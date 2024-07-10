package main

type GreetWithBodyRequest struct {
	Text        string `json:"text" validate:"required,min=3,max=100,textValidator"`
	RepeatTimes int    `json:"repeatTimes" validate:"required,gte=1,lte=5"`
}

type GreetWithQueryRequest struct {
	Text        string `query:"text" validate:"required,min=3,max=100,queryTextValidator"`
	RepeatTimes int    `query:"repeatTimes" validate:"required,gte=1,lte=5"`
}

type GreetWithParamRequest struct {
	Text        string `param:"text" validate:"required,min=3,max=100,queryTextValidator"`
	RepeatTimes int    `query:"repeatTimes" validate:"required,gte=1,lte=5"`
}
