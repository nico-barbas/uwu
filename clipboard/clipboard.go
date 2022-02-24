package clipboard

// NOTE: Only supports Windows for now

func ReadClipboard() (string, error) {
	return readClipboard()
}
