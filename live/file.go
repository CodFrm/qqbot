package live

import "os"

func ShowSource() []string {
	d, err := os.ReadDir("./data/live/source")
	if err != nil {
		return []string{}
	}
	ret := make([]string, 0)
	for _, v := range d {
		ret = append(ret, v.Name())
	}
	return ret
}

func ShowLive() []string {
	d, err := os.ReadDir("./data/live/flv")
	if err != nil {
		return []string{}
	}
	ret := make([]string, 0)
	for _, v := range d {
		ret = append(ret, v.Name())
	}
	return ret
}
