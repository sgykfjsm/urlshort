package urlshort

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"
)

type PathAndURL struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if realPath, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, realPath, 301)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYaml(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)

	return MapHandler(pathMap, fallback), nil
}

func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJson, err := parseJson(jsn)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedJson)

	return MapHandler(pathMap, fallback), nil
}

func buildMap(urls []PathAndURL) map[string]string {
	pathMap := make(map[string]string, len(urls))
	for _, u := range urls {
		pathMap[u.Path] = u.URL
	}

	return pathMap
}

func parseYaml(yml []byte) ([]PathAndURL, error) {
	var pathAndURLs []PathAndURL
	if err := yaml.Unmarshal(yml, &pathAndURLs); err != nil {
		return nil, err
	}

	return pathAndURLs, nil
}

func parseJson(jsn []byte) ([]PathAndURL, error) {
	var pathAndURLs []PathAndURL
	if err := json.Unmarshal(jsn, &pathAndURLs); err != nil {
		return nil, err
	}

	return pathAndURLs, nil
}
