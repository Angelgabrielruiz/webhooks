package domain

type PullRequestEventPayload struct {
	Action      string      `json:"action"`
	PullRequest PullRequest `json:"pull_request"`
	Repository  Repository  `json:"repository"`
	Base        Branch      `json:"base"`
}

type PullRequest struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Head  Branch `json:"head"`
	Base  Branch `json:"base"`
	URL   string `json:"url"`
	User  User   `json:"user"`
}

type Branch struct {
	Ref string `json:"ref"`
	Sha string `json:"sha"`
}

type Repository struct {
	FullName string `json:"full_name"`
	URL      string `json:"url"`
}

type User struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
	Type  string `json:"type"`
	URL   string `json:"url"`
}
