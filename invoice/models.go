package invoice

type Play struct {
	Name     string
	PlayType string
}

type Performance struct {
	PlayID   string
	Audience int
}

type Invoice struct {
	Customer     string
	Performances []Performance
}
