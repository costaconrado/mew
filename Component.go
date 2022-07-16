package mew

import (
	"fmt"
	"reflect"
	"unsafe"
)

var componentsId = map[reflect.Type]uint{}
var maskComponents = map[Mask]reflect.Type{}
var systems = map[reflect.Type]MemSystem{}
var componentsLastId uint

// AddComponent adds any Component to the specified Entity. Component is a simple data struct.
func AddComponent[T any](entity Entity) *T {
	tcomp := typeOf[T]()
	cid := getCompId(tcomp)

	entities[entity] = entities[entity] | 1<<cid

	if systems[tcomp] == nil {

		systems[tcomp] = NewMemSystem(ID(cid), reflect.New(tcomp).Interface(), 1)
		maskComponents[1<<cid] = tcomp
	}
	point := systems[tcomp].New(entity)

	updateFilters(true, EntityMaskPair{entity, entities[entity]})
	return (*T)(point)
}

func getCompId(comp reflect.Type) uint {
	tcomp := reflect.TypeOf(comp)
	if _, ok := componentsId[tcomp]; !ok {
		componentsLastId++
		componentsId[tcomp] = componentsLastId
		// if systems[tcomp] == nil {
		// 	systems[tcomp] = NewMemSystem(ID(componentsId[tcomp]), comp, 1)
		// 	maskComponents[1<<componentsId[tcomp]] = tcomp
		// }
	}
	return componentsId[tcomp]
}

func DelComponent[T any](entity Entity) {
	tcomp := typeOf[T]()
	systems[tcomp].Recycle(entity)
	entities[entity] &^= 1 << componentsId[tcomp]
	updateFilters(false, EntityMaskPair{entity, entities[entity]})
}

func remComponentsFromEntities(ents ...Entity) {
	for _, entity := range ents {
		entmask := entities[entity]
		for mask, tcomp := range maskComponents {
			if mask&entmask != 0 {
				systems[tcomp].Recycle(entity)
			}
		}
		entities[entity] = 0
	}
}

func GetComponent[T any](entity Entity) unsafe.Pointer {
	tcomp := typeOf[T]()
	if systems[tcomp] == nil {
		panic(fmt.Sprintf("Component [%s] does not exist in memory pool", tcomp.Name()))
	}
	return systems[tcomp].Get(entity)
}

func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}
