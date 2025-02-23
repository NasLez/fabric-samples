package domain

type EblSelector struct {
	Selector  map[string]string `json:"selector"`
	UserIndex []string          `json:"user_index"`
}
