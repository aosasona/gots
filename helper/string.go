package helper

func WithDefaultString(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}

	return value
}
