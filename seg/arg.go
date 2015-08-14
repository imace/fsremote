package seg

type SegoArg struct {
	Text      string
	IsNSearch bool
}

type SegoToken struct {
	Text      string      `json:"text"`
	Frequency int         `json:"frequency"`
	Pos       string      `json:"pos"`
	Tokens    []SegoToken `json:"tokens,omitempty"`
}
