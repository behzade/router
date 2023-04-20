package router

type Param struct {
	key   string
	value string
}

type Params []Param

func (p *Params) Get(key string) string {
	for _, v := range *p {
		if v.key == key {
			return v.value
		}
	}

	return ""
}
