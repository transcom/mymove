package auth

// Feature is a gated feature that's accessible to some users
type Feature string

const (
	// FeatureDPS allows access to DPS authentication related features
	FeatureDPS Feature = "dps"
)
