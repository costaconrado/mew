package mew

import (
	"fmt"
	"reflect"
	"sort"
)

type Mask uint64
type Entities []Entity

func (e Entities) Len() int           { return len(e) }
func (e Entities) Less(a, b int) bool { return e[a] < e[b] }
func (e Entities) Swap(a, b int)      { e[a], e[b] = e[b], e[a] }

var filters = []*entityFilter{}

type entityFilter struct {
	mask     Mask
	entities Entities
}

type Filter interface {
	Entities() []Entity
}

type EntityMaskPair struct {
	entity ID
	mask   Mask
}

func (e *entityFilter) Entities() []Entity {
	return e.entities
}

func (e *entityFilter) sort() {
	sort.Sort(e.entities)
}

func (e *entityFilter) indexOf(entity Entity, limit int) int {
	index := sort.Search(limit, func(ind int) bool { return e.entities[ind] >= entity })
	if index < limit && e.entities[index] == entity {
		return index
	}
	return limit
}

func NewFilter(components ...interface{}) Filter {
	mask := Mask(0)
	for _, comp := range components {
		mask |= 1 << getCompId(reflect.TypeOf(comp))
		fmt.Printf("adding compoent [%s] to filter. New masl [%06b]\n", reflect.TypeOf(comp).Name(), mask)
	}

	filter := &entityFilter{
		mask:     mask,
		entities: make([]Entity, 0),
	}

	collectEntities(filter)

	filters = append(filters, filter)
	return filter
}

func collectEntities(filter *entityFilter) {
	for ent, mask := range entities {
		if filter.mask&mask == filter.mask {
			filter.entities = append(filter.entities, ent)
		}
	}
	filter.sort()
}

func updateFilters(add bool, entityPairs ...EntityMaskPair) {
	if add {
		for _, filter := range filters {
			if batchInsertToFilter(filter, entityPairs...) {
				filter.sort()
			}
		}
	} else {
		for _, filter := range filters {
			if batchRemoveFromFilter(filter, entityPairs...) {
				filter.sort()
			}
		}
	}
}

func batchRemoveFromFilter(filter *entityFilter, entityPairs ...EntityMaskPair) (needSort bool) {
	var entityRemoveIndex = []int{}
	needSort = false

	for _, entityPair := range entityPairs {
		rmask := entityPair.mask
		emask := entities[entityPair.entity]
		newmask := emask &^ rmask

		if filter.mask&newmask != filter.mask {
			index := filter.indexOf(entityPair.entity, len(filter.entities))
			entityRemoveIndex = append(entityRemoveIndex, index)
		}
	}

	if len(entityRemoveIndex) > 0 {
		for i, index := range entityRemoveIndex {
			filter.entities[index] = filter.entities[len(filter.entities)-1-i]
		}
		filter.entities = filter.entities[:len(filter.entities)-len(entityRemoveIndex)]
		needSort = true
	}

	return
}

func batchInsertToFilter(filter *entityFilter, entityPairs ...EntityMaskPair) (needSort bool) {
	needSort = false
	filterEntityCount := len(filter.entities)

	for _, entityPair := range entityPairs {
		amask := entityPair.mask
		emask := entities[entityPair.entity]
		newmask := emask | amask
		entity := entityPair.entity

		if filter.mask&newmask == filter.mask {
			index := filter.indexOf(entityPair.entity, filterEntityCount)
			if index == filterEntityCount {
				filter.entities = append(filter.entities, entity)
				needSort = true
			}
		}
	}
	return
}
