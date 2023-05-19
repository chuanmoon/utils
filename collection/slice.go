package collection

func SliceReverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func SliceChunk[T any](slice []T, chunkLen int) [][]T {
	var chunks [][]T
	lenArr := len(slice)
	for i := 0; i < lenArr; i += chunkLen {
		end := i + chunkLen
		if end > lenArr {
			end = lenArr
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}
