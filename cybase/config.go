package cybase

func ReadConfig(key, defaultValue string) string {
	if len(matedata) == 0 {
		panic("matedata is empty, please call cybase.Init() first")
	}
	value, ok := matedata[key]
	if !ok {
		return defaultValue
	}
	return value
}
