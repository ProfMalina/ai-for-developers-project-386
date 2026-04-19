package app

type StorageMode string

const (
	StorageModePostgres StorageMode = "postgres"
	StorageModeMemory   StorageMode = "memory"
)
