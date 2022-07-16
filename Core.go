package mew

var isInitialized bool

// Init initializes all of your ECS updateSystems. Call it once on game world start before calling Update
func Init() {
	if isInitialized {
		panic("ECS already initialized!")
	}

	initLayers()
	isInitialized = true
}

// Update calls Update of each UpdateSystem every frame. Used to handle most of game logic.
func Update() {
	updateLayers()
}
