package raml

import (
	"errors"
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

var header = `#%RAML 1.0
---
`

type RAML struct {
	Title         string          `yaml:"title,omitempty"`
	BaseUri       string          `yaml:"baseUri,omitempty"`
	Protocols     []string        `yaml:"protocols,omitempty"`
	MediaType     string          `yaml:"mediaType,omitempty"`
	Version       string          `yaml:"version,omitempty"`
	Documentation []Documentation `yaml:"documentation,omitempty"`

	Resources `yaml:",inline"`
}

func (r *RAML) String() string {
	bytes, _ := yaml.Marshal(r)
	return fmt.Sprintf("%s%s", header, bytes)
}

type Documentation struct {
	Title   string `yaml:"title"`
	Content string `yaml:"content"`
}

type Resources map[string]*Resource

type Resource struct {
	DisplayName     string    `yaml:"displayName,omitempty"`
	Description     string    `yaml:"description,omitempty"`
	Responses       Responses `yaml:"responses,omitempty"`
	Body            Body      `yaml:"body,omitempty"`
	Is              []string  `yaml:"is,omitempty"`
	Type            string    `yaml:"type,omitempty"`
	SecuredBy       []string  `yaml:"securedBy,omitempty"`
	UriParameters   []string  `yaml:"uirParameters,omitempty"`
	QueryParameters []string  `yaml:"queryParameters,omitempty"`

	Resources `yaml:",inline"`
}

type Responses map[int]Response

type Response struct {
	Body `yaml:"body,omitempty"`
}

type Body map[string]Example // Content-Type to Example

type Example struct {
	Example string `yaml:"example,omitempty"`
}

func (r *RAML) Add(method string, route string, resource *Resource) error {
	if resource == nil {
		return errors.New("raml.Add(): resource can't be nil")
	}
	if r.Resources == nil {
		r.Resources = Resources{}
	}

	return r.Resources.upsert(method, route, resource)
}

func (r *RAML) AddUnder(parentRoute string, method string, route string, resource *Resource) error {
	if resource == nil {
		return errors.New("raml.Add(): resource can't be nil")
	}
	if r.Resources == nil {
		r.Resources = Resources{}
	}

	if parentRoute == "" || parentRoute == "/" {
		return errors.New("raml.AddUnderParent(): parentRoute can't be empty or '/'")
	}

	if !strings.HasPrefix(route, parentRoute) {
		return errors.New("raml.AddUnderParent(): parentRoute must be present in the route string")
	}

	route = strings.TrimPrefix(route, parentRoute)
	if route == "" {
		route = "/"
	}

	parentNode, found := r.Resources[parentRoute]
	if !found {
		parentNode = &Resource{
			Resources: Resources{},
			Responses: Responses{},
		}
		r.Resources[parentRoute] = parentNode
	}

	return parentNode.Resources.upsert(method, route, resource)
}

// Find or create node tree from a given route and inject the resource.
func (r Resources) upsert(method string, route string, resource *Resource) error {
	currentNode := r

	parts := strings.Split(route, "/")
	if len(parts) > 0 {
		last := len(parts) - 1

		// Upsert route of the resource.
		for _, part := range parts[:last] {
			if part == "" {
				continue
			}
			part = "/" + part

			node, found := currentNode[part]
			if !found {
				node = &Resource{
					Resources: Resources{},
					Responses: Responses{},
				}

				currentNode[part] = node
			}
			currentNode = node.Resources
		}

		if parts[last] != "" {
			// Upsert resource into the very bottom of the node tree.
			part := "/" + parts[last]
			node, found := currentNode[part]
			if !found {
				node = &Resource{
					Resources: Resources{},
					Responses: Responses{},
				}
			}
			currentNode[part] = node
			currentNode = node.Resources
		}
	}

	method = strings.ToLower(method)
	if _, found := currentNode[method]; found {
		return nil
		// return fmt.Errorf("duplicated method route: %v %v", method, route)
	}

	currentNode[method] = resource

	return nil
}
