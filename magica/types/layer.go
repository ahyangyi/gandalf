package types

type Layer struct {
	LayerID    int
	Attributes Dictionary
}

func (r *MagicaReader) GetLayer() Layer {
	layerID := r.GetInt32()
	attributes := r.GetDictionary()

	return Layer{
		LayerID:    layerID,
		Attributes: attributes,
	}
}

func (s *Layer) IsChunk() bool {
	return true
}

func (s *Layer) GetChunkName() string {
	return "LAYR"
}
