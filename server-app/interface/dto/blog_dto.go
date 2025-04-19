package dto

type BlogPost struct {
	ID      string `json:"id"`
	UserID  string `json:"userId" binding:"required,min=2,max=10"`
	Title   string `json:"title" binding:"required,min=1,max=50"`
	Content string `json:"content" binding:"required,min=1,max=8000"`
}

type BlogPostResponse struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
	Created string `json:"created_at"`
}

type BlogCreatedResponse struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}
