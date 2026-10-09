package url_cache

func linkKey(shortCode string) string {
	return "link:" + shortCode
}

func counterKey(userID string) string {
	return "counter:" + userID
}

func limitKey(userID string) string {
	return "limit:" + userID
}
