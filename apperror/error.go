package apperror

import (
	"reflect"
)

type Error struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
	Type    string   `json:"type"`
	Context Context  `json:"context"`
}

type Context struct {
	Peers           []string  `json:"peers,omitempty"`
	PeersWithLabels []string  `json:"peersWithLabels,omitempty"`
	Label           string    `json:"label"`
	Key             string    `json:"key,omitempty"`
	Value           *struct{} `json:"value,omitempty"`
}

func New(msg, ty string, path string) Error {
	return Error{
		Message: msg,
		Path:    []string{path},
		Type:    ty,
		Context: Context{
			Label: path,
			Key:   path,
		},
	}
}

func NewWithPeers(msg, ty string, path []string, label string, peers []string) Error {
	return Error{
		Message: msg,
		Type:    ty,
		Path:    path,
		Context: Context{
			Peers:           peers,
			PeersWithLabels: peers,
			Label:           label,
			Value:           new(struct{}),
		},
	}
}

func (e Error) Error() string {
	return e.Message
}

func FromError(err error) (bool, Error) {
	val := reflect.ValueOf(err)
	if val.Type().String() != "apperror.Error" {
		return false, Error{}
	}

	iface := val.Interface()

	return true, iface.(Error)
}
