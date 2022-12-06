package dto

type CarouselItem struct {
	PosterImage     string `json:"poster_image"`
	BackgroundImage string `json:"background_image"`
	Title           string `json:"title"`
	Content         string `json:"content"`
	Value           string `json:"value"`
}

type HomepageContentResponse struct {
	CarouselItems []*CarouselItem `json:"carousel_items"`
}
