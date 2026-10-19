package urlbuilder

import "net/url"

func Build(base string, path string, queryParams map[string]string) (string, error) {
	baseURL, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	joined := baseURL.JoinPath(path).String()
	query := url.Values{}
	for key, value := range queryParams {
		query.Add(key, value)
	}
	if query.Encode() != "" {
		return joined + "?" + query.Encode(), nil
	}
	return joined, nil
}
