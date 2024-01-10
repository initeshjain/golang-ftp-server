package models

// Config holds configuration parameters
type Config struct {
	UploadDir string `json:"UploadDir"`
	DbPath    string `json:"DbPath"`
}

type UploadResponse struct {
	Url                string `json:"url"`
	Pathname           string `json:"pathname"`
	ContentType        string `json:"contentType"`
	ContentDisposition string `json:"contentDisposition"`
	Message            string `json:"message"`
}

type DeleteResponse struct {
	Message string `json:"message"`
}

type GetResponse struct {
}
