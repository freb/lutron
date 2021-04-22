package lutron

import "encoding/json"

func pretty(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

func prettyb(b []byte) string {
	v := make(map[string]interface{})
	json.Unmarshal(b, &v)
	return pretty(v)
}
