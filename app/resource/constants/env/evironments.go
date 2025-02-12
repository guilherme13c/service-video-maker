package env

type Environment string

const (
	Development Environment = "DEVELOPMENT"
	Production  Environment = "PRODUCTION"
	STAGE       Environment = "STAGE"
)
