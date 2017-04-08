package util

import (
	"strconv"
	"strings"
)

func Id2Path(id int64) string {

	parts := []string{}
	input := strconv.FormatInt(id, 10)

	for len(input) > 3 {

		chunk := input[0:3]
		input = input[3:]
		parts = append(parts, chunk)
	}

	if len(input) > 0 {
		parts = append(parts, input)
	}

	path := strings.Join(parts, "/")
	return path
}
