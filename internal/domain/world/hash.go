package world

import (
	"encoding/binary"
	"hash/fnv"
)

// hashCoordinates generates a deterministic hash from world seed and tile coordinates
// Used to seed RNG for consistent world generation across runs
func hashCoordinates(worldSeed int64, chunkX, chunkY, localTileX, localTileY int) int64 {
	h := fnv.New64a()

	// Write all inputs to hash in a deterministic order
	binary.Write(h, binary.LittleEndian, worldSeed)
	binary.Write(h, binary.LittleEndian, int64(chunkX))
	binary.Write(h, binary.LittleEndian, int64(chunkY))
	binary.Write(h, binary.LittleEndian, int64(localTileX))
	binary.Write(h, binary.LittleEndian, int64(localTileY))

	return int64(h.Sum64())
}
