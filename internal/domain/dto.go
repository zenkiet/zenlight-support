package domain

type InstallFileDTO struct {
	Name      string `json:"name"`
	Data      []byte `json:"data"`
	Extension string `json:"extension"`
}
