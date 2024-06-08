package model

type Err struct {
	Code    int
	Message string
}

func (e *Err) Error() string {
	return e.Message
}

type Genre struct {
	Id   string
	Name string
}

type Author struct {
	Id   string
	Name string
}

type Chapter struct {
	Id      string
	Title   string
	Content string
}

type Novel struct {
	Id            string
	Title         string
	Rate          float32
	Author        []*Author
	Genre         []*Genre
	CoverImage    string
	Description   string
	Status        string
	Chapters      []*Chapter
	LatestChapter string
}
