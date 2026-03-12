package app

type TrainingArtifact struct {
	FileName     string `json:"fileName"`
	DownloadName string `json:"downloadName"`
	Content      string `json:"content"`
	Focus        string `json:"focus"`
}

type Question struct {
	ID             string `json:"id"`
	Prompt         string `json:"prompt"`
	ExpectedAnswer string `json:"expectedAnswer"`
}

type TrainingSet struct {
	ID            string           `json:"id"`
	Title         string           `json:"title"`
	Description   string           `json:"description"`
	Difficulty    string           `json:"difficulty"`
	CreatedAt     string           `json:"createdAt"`
	QuestionCount int              `json:"questionCount"`
	Tags          []string         `json:"tags"`
	Artifact      TrainingArtifact `json:"artifact"`
	Questions     []Question       `json:"questions"`
}

type QuestionScoreResult struct {
	DatasetID       string `json:"datasetId"`
	DatasetTitle    string `json:"datasetTitle"`
	QuestionID      string `json:"questionId"`
	Prompt          string `json:"prompt"`
	SubmittedAt     string `json:"submittedAt"`
	SubmittedAnswer string `json:"submittedAnswer"`
	ExpectedAnswer  string `json:"expectedAnswer"`
	Correct         bool   `json:"correct"`
}

type MistakeDetail struct {
	QuestionID      string `json:"questionId"`
	Prompt          string `json:"prompt"`
	ExpectedAnswer  string `json:"expectedAnswer"`
	SubmittedAnswer string `json:"submittedAnswer"`
}

type ScoreResult struct {
	DatasetID    string          `json:"datasetId"`
	DatasetTitle string          `json:"datasetTitle"`
	SubmittedAt  string          `json:"submittedAt"`
	Score        int             `json:"score"`
	Total        int             `json:"total"`
	Passed       bool            `json:"passed"`
	CorrectCount int             `json:"correctCount"`
	WrongCount   int             `json:"wrongCount"`
	Mistakes     []MistakeDetail `json:"mistakes"`
}

type CreateTrainingInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty"`
}
