package storage

func Disk(dir string) IPrivateStorage {
	if dir == "tmp" || dir == "temp" {
		return Temp()
	}

	return newStorage(dir)
}

func Temp() IPrivateStorage {
	return newTempStorage()
}
