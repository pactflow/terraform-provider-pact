package broker

type Pacticipant struct {
	Name          string `json:"name,omitempty"`
	RepositoryURL string `json:"repositoryUrl,omitempty"`
}
