package nat64

type Nat64Stats struct {
	translationsActive    float64 `default:"0"`
	translationsExpired   float64 `default:"0"`
	sessionsFound         float64 `default:"0"`
	sessionsCreated       float64 `default:"0"`
	packetsTranslated4to6 float64 `default:"0"`
	packetsTranslated6to4 float64 `default:"0"`
}
