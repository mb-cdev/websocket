package websocket

func cutByteSlice(data []byte, n int) [][]byte {
	r := make([][]byte, 0)

	if len(data) == 0 {
		return r
	}

	if n >= len(data) {
		r = append(r, data)
		return r
	}

	var offset int
	for {
		//buff := make([]byte, 0)
		endCutIndex := offset + n

		if endCutIndex >= len(data) {
			r = append(r, data[offset:])
			break
		}

		r = append(r, data[offset:endCutIndex])

		offset += n
	}

	return r
}
