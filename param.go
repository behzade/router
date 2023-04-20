package router

type Param struct {
	key   string
	value []byte
}

type Params []Param

func (p *Params) Get(key string) []byte {
	for _, v := range *p {
		if v.key == key {
			return v.value
		}
	}

	return nil
}
