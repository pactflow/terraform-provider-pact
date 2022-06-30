package broker

type Pacticipant struct {
	Name          string `json:"name,omitempty" pact:"example=terraform-client"`
	RepositoryURL string `json:"repositoryUrl,omitempty" pact:"example=https://github.com/pactflow/terraform-provider-pact"`
	MainBranch    string `json:"mainBranch,omitempty" pact:"example=main"`
	DisplayName   string `json:"displayName,omitempty" pact:"example=terraform client"`
}
