package models

type Level string

const (
	BRONZE  Level = "BRONZE"
	SILVER  Level = "SILVER"
	GOLD    Level = "GOLD"
	UNKNOWN Level = ""
)

func GetLevelFromName(name string) Level {
	switch name {
	case "BRONZE":
		return BRONZE
	case "SILVER":
		return SILVER
	case "GOLD":
		return GOLD
	}

	return UNKNOWN
}
