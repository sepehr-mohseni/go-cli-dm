package manager

import "net/http"

func (dm Download) GetNewRequest(method string) (*http.Request, error) {
	r, err := http.NewRequest(
		method,
		dm.Url,
		nil,
	)

	if err != nil {
		return nil, err
	}

	r.Header.Set("User-Agent", "Sep DM V1")
	return r, nil
}