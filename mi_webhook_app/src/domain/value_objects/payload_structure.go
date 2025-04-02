// File: src/domain/value_objects/payload_structure.go
package domain

import "time" // Importa time

// --- Pull Request Event ---

type PullRequestEventPayload struct {
	Action      string      `json:"action"`
	Number      int         `json:"number"`
	PullRequest PullRequest `json:"pull_request"`
	Repository  Repository  `json:"repository"`
	Sender      User        `json:"sender"`
}

type PullRequest struct {
	ID          int         `json:"id"`
	HTMLURL     string      `json:"html_url"` // URL para el navegador
	Title       string      `json:"title"`
	User        User        `json:"user"` // Quién creó el PR
	State       string      `json:"state"`
	Merged      bool        `json:"merged"`
	MergedAt    *time.Time  `json:"merged_at"` // Puntero si puede ser null
	CreatedAt   *time.Time  `json:"created_at"` // Puntero si puede ser null
	UpdatedAt   *time.Time  `json:"updated_at"` // Puntero si puede ser null
	Head        Branch      `json:"head"`
	Base        Branch      `json:"base"`
}

type Branch struct {
	Ref  string     `json:"ref"`
	Sha  string     `json:"sha"`
	Repo Repository `json:"repo"`
}

type Repository struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"` // owner/repo
	HTMLURL  string `json:"html_url"`  // URL al repo
}

type User struct {
	Login   string `json:"login"`
	ID      int    `json:"id"`
	HTMLURL string `json:"html_url"` // URL al perfil
	Type    string `json:"type"`
}

// --- Workflow Run Event ---

type WorkflowRunEventPayload struct {
	Action      string      `json:"action"` // "requested" o "completed"
	WorkflowRun WorkflowRun `json:"workflow_run"`
	Workflow    Workflow    `json:"workflow"`
	Repository  Repository  `json:"repository"`
	Sender      User        `json:"sender"`
}

type WorkflowRun struct {
	ID           int64               `json:"id"`
	Name         string              `json:"name"`
	HeadBranch   string              `json:"head_branch"`
	HeadSha      string              `json:"head_sha"`
	RunNumber    int                 `json:"run_number"`
	Event        string              `json:"event"`
	Status       string              `json:"status"`
	Conclusion   string              `json:"conclusion"` // Puede ser null si no está completed
	WorkflowID   int64               `json:"workflow_id"`
	HTMLURL      string              `json:"html_url"` // URL a la ejecución específica
	CreatedAt    time.Time           `json:"created_at"` // GitHub usualmente lo envía no-null
	UpdatedAt    time.Time           `json:"updated_at"` // GitHub usualmente lo envía no-null
	PullRequests []WorkflowPullRequest `json:"pull_requests"`
}

type Workflow struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Path  string `json:"path"`
	State string `json:"state"`
}

type WorkflowPullRequest struct {
	URL    string `json:"url"`
	ID     int64  `json:"id"`
	Number int    `json:"number"`
	Head   struct {
		Ref  string `json:"ref"`
		Sha  string `json:"sha"`
		Repo struct {
			ID   int64  `json:"id"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"repo"`
	} `json:"head"`
	Base struct {
		Ref  string `json:"ref"`
		Sha  string `json:"sha"`
		Repo struct {
			ID   int64  `json:"id"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"repo"`
	} `json:"base"`
}