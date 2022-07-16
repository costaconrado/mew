package mew

type Entity = ID

var recycleIDs []ID = []ID{}
var entities map[Entity]Mask = map[Entity]Mask{}

// type Entity uint32

var currentEntityId Entity

// NewEntity creates new Entity in the game world.
func NewEntity() Entity {
	var id ID

	if len(recycleIDs) > 0 {
		id = recycleIDs[len(recycleIDs)-1]
		recycleIDs = recycleIDs[:len(recycleIDs)-1]
	} else {
		currentEntityId++
		id = currentEntityId
	}
	entities[id] = 0

	return id
}

func DeleteEntity(entity Entity) {
	if entity < 1 {
		LogMessage("[World.RemEntity] invalid entity id %d\n", entity)
		return
	}
	remComponentsFromEntities(entity)
	recycleEntitiesAndUpdateFilters(entity)
}

func recycleEntitiesAndUpdateFilters(ents ...Entity) {
	for _, entity := range ents {
		updateFilters(false, EntityMaskPair{entity, entities[entity]})
		entities[entity] = Mask(0)
	}
	recycleIDs = append(recycleIDs, ents...)
}
