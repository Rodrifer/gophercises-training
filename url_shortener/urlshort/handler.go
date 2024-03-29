package urlshort

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path, ok := pathsToUrls[r.URL.Path]
		if ok {
			http.Redirect(w, r, path, http.StatusFound)
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
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

//
// JSONHandler implementation
//
func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJSON, err := parseJSON(jsn)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	pathMap := buildMap(parsedJSON)
	return MapHandler(pathMap, fallback), nil
}

//
// BoltDBHandler handler
//
func BoltDBHandler(pathMap map[string]string, fallback http.Handler) (http.HandlerFunc, error) {
	return MapHandler(pathMap, fallback), nil
}

func parseYAML(yml []byte) (prsd []map[string]string, err error) {
	err = yaml.Unmarshal(yml, &prsd)
	return prsd, err
}

func parseJSON(jsn []byte) (prsd []map[string]string, err error) {
	err = json.Unmarshal(jsn, &prsd)
	return prsd, err
}

func buildMap(parsedyml []map[string]string) map[string]string {
	mergedMap := make(map[string]string)
	for _, entry := range parsedyml {
		key := entry["path"]
		mergedMap[key] = entry["url"]
	}
	return mergedMap
}
