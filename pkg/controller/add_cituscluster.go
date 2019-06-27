package controller

import (
	"github.com/infinivision/citus-operator/pkg/controller/cituscluster"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, cituscluster.Add)
}
