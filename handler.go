package urlshort

import (
	"encoding/json"
	"net/http"

	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v3"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusPermanentRedirect)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

type pathInfo []struct {
	Path string `yaml:"path" json:"path"`
	Url  string `yaml:"url" json:"url"`
}

func parseYAML(yamlData []byte) (pathInfo, error) {
	var p pathInfo
	err := yaml.Unmarshal(yamlData, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func parseJSON(jsonData []byte) (pathInfo, error) {
	var p pathInfo
	err := json.Unmarshal(jsonData, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func pathsToMap(paths pathInfo) map[string]string {
	m := make(map[string]string)
	for _, path := range paths {
		m[path.Path] = path.Url
	}
	return m
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
func YAMLHandler(yamlData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	paths, err := parseYAML(yamlData)
	if err != nil {
		return nil, err
	}
	m := pathsToMap(paths)
	return MapHandler(m, fallback), nil
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
//		[
//			{
//	    		"path": "/some-path"
//	    		"url": "https://www.some-url.com/demo"
//			}
//		]
//
// The only errors that can be returned all related to having
// invalid JSON data.
func JSONHandler(jsonData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	paths, err := parseJSON(jsonData)
	if err != nil {
		return nil, err
	}
	m := pathsToMap(paths)
	return MapHandler(m, fallback), nil
}

// DBHandler will use the provided BoltDB and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the DB, then the
// fallback http.Handler will be called instead.
// There is a test file TestDB.db, which stores test data:
//
//	/bolt : https://github.com/boltdb/bolt#opening-a-database
//	/bolt-godoc : https://pkg.go.dev/github.com/boltdb/bolt#Open
func DBHandler(db *bolt.DB, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var path string
		err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("paths"))
			val := b.Get([]byte(r.URL.Path))
			if val != nil {
				path = string(val)
			}
			return nil
		})
		if err == nil && path != "" {
			http.Redirect(w, r, path, http.StatusPermanentRedirect)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}
