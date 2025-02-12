package status

type Status string

const (
	Pending    Status = "PENDING"
	Processing Status = "PROCESSING"
	Completed  Status = "COMPLETED"
	Error      Status = "ERROR"
)
