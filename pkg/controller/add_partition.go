package controller

import (
	"github.com/adam-cattermole/striot-operator/pkg/controller/partition"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, partition.Add)
}
