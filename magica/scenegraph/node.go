package scenegraph

import (
	"github.com/ahyangyi/gandalf/geometry"
	"github.com/ahyangyi/gandalf/magica/types"
)

type Map map[int]types.SceneGraphItem

type Model struct {
	Points types.PointData
	Size   types.Size
}

type Node struct {
	Location geometry.Point
	Size     types.Size
	Models   []Model
	Children []Node
}

func GetScenegraph(scenegraphMap Map, allowedLayers []int, disallowedLayerNames []string, pointData []types.PointData, sizeData []types.Size) Node {
	if len(scenegraphMap) == 0 && len(sizeData) > 0 && len(pointData) > 0 {
		return Node{
			Location: geometry.Point{},
			Size:     sizeData[0],
			Models:   []Model{{Points: pointData[0], Size: sizeData[0]}},
		}
	} else if len(scenegraphMap) > 0 {
		return Compose(scenegraphMap, scenegraphMap[0], 0, 0, 0, -1, allowedLayers, "", disallowedLayerNames, pointData, sizeData)
	}

	return Node{}
}

func (n *Node) GetExtents() Extent {
	extents := make(Extents, len(n.Models))
	for idx, model := range n.Models {
		extents[idx] = Extent{
			Min: geometry.Point{X: n.Location.X, Y: n.Location.Y, Z: n.Location.Z},
			Max: geometry.Point{X: n.Location.X + model.Size.X, Y: n.Location.Y + model.Size.Y, Z: n.Location.Z + model.Size.Z},
		}
	}

	for _, child := range n.Children {
		extents = append(extents, child.GetExtents())
	}

	return extents.GetBounds()
}
