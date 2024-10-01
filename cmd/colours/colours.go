package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

func floatToHex(inp float64) string {
	return fmt.Sprintf("%02X", int(math.Round(inp*255)))
}

func colourRename(inp string) (string, string) {
	split := strings.SplitN(inp, "/", 2)
	colourName := split[0]
	colourVal := split[1]
	if colourVal == "1000" {
		colourVal = "950"
	}
	return strings.ToLower(colourName), colourVal
}

func main() {
	defer os.Stdin.Close()
	input, _ := io.ReadAll(os.Stdin)

	data := make([]map[string]any, 0)
	_ = json.Unmarshal(input, &data)
	outputMap := make(map[string]map[string]string, 0)
	for _, v := range data {
		colourMap := v["rgba"].(map[string]any)
		name, value := colourRename(v["name"].(string))
		code := fmt.Sprintf(
			"#%s%s%s",
			floatToHex(colourMap["r"].(float64)),
			floatToHex(colourMap["g"].(float64)),
			floatToHex(colourMap["b"].(float64)),
		)
		_, ok := outputMap[name]
		if !ok {
			outputMap[name] = make(map[string]string, 0)
		}

		outputMap[name][value] = code
	}

	outBytes, _ := json.Marshal(outputMap)
	os.Stdout.Write(outBytes)

}
