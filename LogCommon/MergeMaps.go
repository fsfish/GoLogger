package loghelper

//将多个map合并为一个map
func MergeMap(param ...map[string]interface{}) map[string]interface{} {
	pCount := len(param)
	if pCount == 0 {
		return nil
	}
	if pCount == 1 {
		return param[0]
	}
	n := make(map[string]interface{}, pCount)
	for i := pCount - 1; i >= 0; i-- {
		p := param[i]
		if p == nil {
			continue
		}
		for k, v := range p {
			if _, ok := n[k]; !ok {
				n[k] = v
			}
		}
	}
	return n
}
