package domain

type Blog struct {
	ID      string
	Title   string
	Content string
	UserID  string
}

type BlogPost struct {
	LoginID string
	Title   string
	Content string
}
