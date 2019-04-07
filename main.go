package main

import 
(
  "io/ioutil"
  "fmt"
  "github.com/paulmach/orb/mercator"
  "github.com/paulmach/go.geojson"
  "github.com/fogleman/gg"
  "time"
)

const (
  width = 1366
  height = 1024
)

func main() {
  t1 := time.Now()
  rawFeatureCollectionJSON, err := ioutil.ReadFile("./Адм-территориальные границы РФ в формате GeoJSON/admin_level_3.geojson")
  if err != nil {
    fmt.Printf("Coulnd't load data.geojson file: %v", err)
    return
  }

  fc := geojson.NewFeatureCollection()
  fc, err = geojson.UnmarshalFeatureCollection(rawFeatureCollectionJSON)
  if err != nil {
    fmt.Printf("Error decoding the data into GeoJSON feature collection: %v", err)
    return
  }

  ctx := gg.NewContext(width, height)
  ctx.SetRGB(0, 0, 0) // Paint it, black!
  ctx.Clear()

  ctx.ClearPath()
  ctx.SetRGB(0.543, 0, 0)
  ctx.SetLineWidth(1.5)

  for _, f := range fc.Features {
    switch {
      case f.Geometry.IsMultiPolygon():
      	fmt.Println("Detected MultiPolygon type of data")
        if ctx, err = DrawMultiPolygon(ctx, f); err != nil {
          fmt.Printf("Couldn't handle MultiPolygon type: %v", err)
          return
        }
      default:
        fmt.Println("Oops, exotic data type")
        return
    }
  }
  err = ctx.SavePNG("image.png")
  if err != nil {
    fmt.Println("Tragic event saving context as PNG: %v", err)
  }

  fmt.Printf("\nProgram finished in %v ", time.Now().Sub(t1))
}

func DrawMultiPolygon(c *gg.Context, f *geojson.Feature) (res *gg.Context, err error) {
	for i := 0; i < len(f.Geometry.MultiPolygon); i++ {
		for j := 0; j < len(f.Geometry.MultiPolygon[i]); j++ {
			val := f.Geometry.MultiPolygon[i][j]
			for k := 0; k < len(val); k++ {
				x, y := mercator.ToPlanar(val[k][0], val[k][1], 10)
				// x := val[k][0] + width / 2
				// y := val[k][1] + height / 2
				c.LineTo(x, y)
			}
			c.FillPreserve() // fills the current path with the current color
		    c.Stroke()
		    c.Fill()
		}
	}

	return c, err
}