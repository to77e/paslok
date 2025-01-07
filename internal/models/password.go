package models

const (
	DefaultLength = 18
	DefaultChunk  = 6
)

type Password struct {
	Length    int
	ChunkSize int
	Uppercase bool
	Special   bool
	Number    bool
	Dash      bool
}

type CreatePasswordRequest struct {
	Service   string `json:"service" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Comment   string `json:"comment" validate:"required"`
	Password  string `json:"password"`
	Length    int    `json:"length" validate:"required,gte=6,lte=64"`
	ChunkSize int    `json:"chunk_size"`
	Uppercase bool   `json:"uppercase"`
	Special   bool   `json:"special"`
	Number    bool   `json:"number"`
	Dash      bool   `json:"dash"`
}

type ReadPasswordRequest struct {
	Id      int64  `json:"id" validate:"required_without=Service,excluded_with=Service"`
	Service string `json:"service" validate:"required_without=Id,excluded_with=Id"`
}

type ListPasswordsRequest struct {
	SearchTerm string `json:"search_term"`
}

type DeletePasswordRequest struct {
	Id      int64  `json:"id" validate:"required_without=Service,excluded_with=Service"`
	Service string `json:"service" validate:"required_without=Id,excluded_with=Id"`
}
