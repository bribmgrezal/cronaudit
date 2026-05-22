package schedule

import "strconv"

// parseInt is a small helper so validate.go stays import-lean.
func parseInt(s string) (int, error) {
	v, err := strconv.Atoi(s)
	return v, err
}
