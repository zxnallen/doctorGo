package pagination

type Page struct {
	Page int `form:"page" json:"page"`
	Size int `form:"size" json:"size"`
}

func (p Page) Normalize() Page {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Size <= 0 {
		p.Size = 10
	}
	if p.Size > 100 {
		p.Size = 100
	}
	return p
}

func (p Page) Offset() int {
	p = p.Normalize()
	return (p.Page - 1) * p.Size
}
