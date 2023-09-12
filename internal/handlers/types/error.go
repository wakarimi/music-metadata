package types

// Error godoc
// @Description Response structure for error messages
// @Property Error (string) Description of the error that occurred
type Error struct {
	Error string `json:"error" binding:"required"`
}
