package structs

type ApiResponse struct {
	Feed struct {
		Entries []Entry `json:"entry"`
	} `json:"feed"`
}

type Entry struct {
	Id struct {
		Attributes struct {
			Id string `json:"im:id"`
		} `json:"attributes"`
	} `json:"id"`

	Name struct {
		Label string `json:"label"`
	} `json:"im:name"`

	Images []Image `json:"im:image"`

	Artist struct {
		Label string `json:"label"`
	} `json:"im:artist"`

	Category struct {
		Attributes struct {
			Label string `json:"label"`
		} `json:"attributes"`
	} `json:"category"`

	Link struct {
		Attributes struct {
			Href string `json:"href"`
		} `json:"attributes"`
	} `json:"link"`

	Price struct {
		Attributes struct {
			Amount string `json:"amount"`
		} `json:"attributes"`
	} `json:"im:price"`

	ReleaseDate struct {
		Label string `json:"label"`
	} `json:"im:releaseDate"`
}

type Image struct {
	Label      string `json:"label"`
	Attributes struct {
		Height string `json:"height"`
	} `json:"attributes"`
}
