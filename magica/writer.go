package magica

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func writeChunkHeader(handle io.Writer, name string, len int) (err error) {
	if _, err = handle.Write([]byte(name)); err != nil {
		return err
	}

	err = binary.Write(handle, binary.LittleEndian, int32(len))
	return
}

func writeChunkHeaderWithChildLength(handle io.Writer, name string, len int, childlen int) (err error) {
	if err = writeChunkHeader(handle, name, len); err != nil {
		return err
	}

	if err = binary.Write(handle, binary.LittleEndian, int32(childlen)); err != nil {
		return err
	}
	return
}

func (v *VoxelObject) writeHeader(handle io.Writer) (err error) {
	err = writeChunkHeader(handle, "VOX ", 150)
	return
}

func (v *VoxelObject) writeMainChunk(handle io.Writer) (err error) {
	mainLen := 52 + len(v.GetPoints())*4 + len(v.PaletteData)
	err = writeChunkHeaderWithChildLength(handle, "MAIN", 0, mainLen)
	return
}

func (v *VoxelObject) writePalette(handle io.Writer) (err error) {
	if err := writeChunkHeaderWithChildLength(handle, "RGBA", len(v.PaletteData), 0); err != nil {
		return err
	}
	_, err = handle.Write(v.PaletteData)
	return
}

func (v *VoxelObject) writeSizeChunk(handle io.Writer) (err error) {
	if err := writeChunkHeaderWithChildLength(handle, "SIZE", 12, 0); err != nil {
		return err
	}

	err = binary.Write(handle, binary.LittleEndian, []int32{int32(v.Size.X), int32(v.Size.Y), int32(v.Size.Z)})
	return
}

func (v *VoxelObject) writeXYZIChunk(handle io.Writer) (err error) {
	points := v.GetPoints()
	if err := writeChunkHeaderWithChildLength(handle, "XYZI", (len(points)*4)+4, 0); err != nil {
		return err
	}

	if err := binary.Write(handle, binary.LittleEndian, int32(len(points))); err != nil {
		return err
	}

	for _, pt := range points {
		if err := binary.Write(handle, binary.LittleEndian, []byte{byte(pt.Point.X), byte(pt.Point.Y), byte(pt.Point.Z), pt.Colour}); err != nil {
			return err
		}
	}

	return
}

func (v *VoxelObject) Save(handle io.Writer) (err error) {
	bw := bufio.NewWriter(handle)
	fns := []func(writer io.Writer) error{v.writeHeader, v.writeMainChunk, v.writeSizeChunk, v.writeXYZIChunk, v.writePalette}
	for _, fn := range fns {
		if err := fn(bw); err != nil {
			return err
		}
	}
	bw.Flush()
	return
}

// SaveToFile saves the voxel object to the specified file
func (v *VoxelObject) SaveToFile(filename string) (err error) {
	handle, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create output file %s: %v", filename, err)
	}

	err = v.Save(handle)
	if err != nil {
		handle.Close()
		return fmt.Errorf("could not open output file: %v", err)
	}

	err = handle.Close()
	return err

}