package anilist

type Response struct {
	Data Data `json:"data"`
}

type Data struct {
	User     User     `json:"User"`
	Activity Activity `json:"Activity"`
	Page     Page     `json:"Page"`
}

type User struct {
	ID int64 `json:"id"`
}

type Activity struct {
	Media     Media `json:"media"`
	CreatedAt int64 `json:"createdAt"`
}

type Media struct {
	CoverImage CoverImage `json:"coverImage"`
	Title      Title      `json:"title"`
	ID         int64      `json:"id"`
}

type CoverImage struct {
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

type Title struct {
	Romaji  string `json:"romaji"`
	English string `json:"english"`
}

type Page struct {
	Activities []Activity `json:"activities"`
}
