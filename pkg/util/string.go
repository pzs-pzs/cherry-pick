package util

func RemoveDuplication(in []string) (out []string) {
	t := make(map[string]struct{})
	for _, v := range in {
		if _, ok := t[v]; ok {
			continue
		}
		out = append(out, v)
		t[v] = struct{}{}
	}
	return
}
