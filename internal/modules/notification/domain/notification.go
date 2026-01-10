package domain

type Headers map[string][]string

type Body struct {
	Type 	string
	Content string
}

func NewHeaders() Headers {
	return make(Headers)
}

func NewBody(bodyType, bodyContent string) *Body {
	return &Body{
		Type: bodyType,
		Content: bodyContent,
	}
}

func (headers Headers) Set(key string, values ...string) {
	headers[key] = values
}

func (headers Headers) Get(key string) []string {
	if values, ok := headers[key]; ok {
		return values
	}

	return []string{}
}

func (headers Headers) GetAtIndex(key string, idx int) string {
	if values, ok := headers[key]; ok {
		return values[idx]
	}

	return ""
}

func (headers Headers) Delete(key string) {
	delete(headers, key)
}

type NotifierService interface {
	Send(Headers, Body) error
}
