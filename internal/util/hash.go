package util

func CheckHashFile(filepath, hash string) bool {
	f := OpenFileMetadata(filepath, true)
	if f.Error != nil {
		return false
	}

	defer DeferCloseFatal(f.File)

	return f.Hash == hash
}
