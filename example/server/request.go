package main

type GreetWithBodyRequest struct {
	Text           string `json:"text" validate:""`
	RepeatTimes    int    `json:"repeatTimes" validate:"required"`
	NecessaryParam string `query:"necessaryParam" validate:"required"`
}

type GreetWithQueryRequest struct {
	Text        string `query:"text" validate:"required,min=3,max=100,queryTextValidator"`
	RepeatTimes int    `query:"repeatTimes" validate:"required,gte=1,lte=5"`
}

type GreetWithParamRequest struct {
	Text        string `param:"text" validate:"required,min=3,max=100,queryTextValidator"`
	RepeatTimes int    `query:"repeatTimes" validate:"required,gte=1,lte=5"`
}
