package mew

import (
	"fmt"
	"reflect"
	"unsafe"
)

var componentsId = map[reflect.Type]uint{}
var maskComponents = map[Mask]reflect.Type{}
var systems map[reflect.Type]MemSystem
var componentsLastId uint

// AddComponent adds any Component to the specified Entity. Component is a simple data struct.
func AddComponent(entity Entity, comp any) {
	if isPointer(comp) {
		panic("Pointers not allowed!")
	}

	tcomp := reflect.TypeOf(comp)
	cid := getCompId(comp)

	entities[entity] = entities[entity] | 1<<cid

	if systems[tcomp] == nil {
		systems[tcomp] = NewMemSystem(entity, comp, 1)
		maskComponents[1<<cid] = tcomp
	}
	systems[tcomp].New(entity)
}

func getCompId(comp any) uint {
	tcomp := reflect.TypeOf(comp)
	if _, ok := componentsId[tcomp]; !ok {
		componentsLastId++
		componentsId[tcomp] = componentsLastId
	}
	return componentsId[tcomp]
}

func DelComponent[T any](entity Entity) {
	tcomp := typeOf[T]()
	systems[tcomp].Recycle(entity)
	entities[entity] &^= 1 << componentsId[tcomp]
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
		panic(fmt.Sprintf("Compent [%s] does not exist in memory pool", tcomp.Name()))
	}
	return systems[tcomp].Get(entity)
}

func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}
