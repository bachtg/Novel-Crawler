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

type GetNovelsRequest struct {
	Page       string
	Keyword    string
	GenreId    string
	CategoryId string
	AuthorId   string
}

type GetNovelsResponse struct {
	Novels  []*Novel
	NumPage int
}

type DetailChapterResponse struct {
	Novel           *Novel
	CurrentChapter  *Chapter
	NextChapter     *Chapter
	PreviousChapter *Chapter
}
