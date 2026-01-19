package math

func (t *TrigTable) Sin(index int) int {
	return t.sin[index%2048]
}

func (t *TrigTable) Cos(index int) int {
	return t.cos[index%2048]
}
