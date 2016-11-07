package source

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Depature struct {
	Line string
	In   int
}

func NextTimes(stop string, place string) ([]Depature, error) {
	reqURL, err := url.Parse("http://widgets.vvo-online.de/abfahrtsmonitor/Abfahrten.do?hst=Joseph%20Stift")
	if err != nil {
		return nil, err
	}

	vals := reqURL.Query()
	vals.Set("hst", stop)
	vals.Set("ort", place)
	reqURL.RawQuery = vals.Encode()

	resp, err := http.Get(reqURL.String())
	if err != nil {
		return nil, err
	}

	var deps [][3]string
	if err := json.NewDecoder(resp.Body).Decode(&deps); err != nil {
		return nil, err
	}

	var ds []Depature
	for _, dep := range deps {
		var in int
		if dep[2] != "" {
			in, err = strconv.Atoi(dep[2])
			if err != nil {
				log.Printf("couldn't parse departure time: %s")
				continue
			}
		}

		ds = append(ds, Depature{
			Line: fmt.Sprintf("%v %v", dep[0], dep[1]),
			In:   in,
		})
	}
	return ds, nil
}
