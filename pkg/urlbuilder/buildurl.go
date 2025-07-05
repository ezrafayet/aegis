package urlbuilder

import "net/url"

func Build(base string, path string) (string, error) {
	baseURL, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	return baseURL.JoinPath(path).String(), nil
}
