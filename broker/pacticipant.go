package broker

type Pacticipant struct {
	Name          string `json:"name,omitempty" pact:"example=terraform-client"`
	RepositoryURL string `json:"repositoryUrl,omitempty" pact:"example=https://github.com/pactflow/terraform-provider-pact"`
	MainBranch    string `json:"mainBranch,omitempty" pact:"example=main"`
	DisplayName   string `json:"displayName,omitempty" pact:"example=terraform client"`
}

// PacticipantsEmbedded contains the embedded pacticipants in the list resource
type PacticipantsEmbedded struct {
	Items []Pacticipant `json:"pacticipants"`
}

// PacticipantsResponse is the response body for the list API call.
// Deliberately doesn't embed HalDoc: the root /pacticipants response has
// array-valued links (e.g. pb:pacticipants), which HalLinks can't unmarshal.
type PacticipantsResponse struct {
	Embedded PacticipantsEmbedded `json:"_embedded"`
}
