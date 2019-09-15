package schema

type TrainingSchema struct {
	Name string            `json:"name"`
	Args map[string]string `json:"args"`
}

type Training struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at"`
}

type TrainingList struct {
	Total int32      `json:"total"`
	Items []Training `json:"items"`
}

type TrainingLogs struct {
	ID         string `json:"id"`
	LogDetails string `json:"log_details"`
}
