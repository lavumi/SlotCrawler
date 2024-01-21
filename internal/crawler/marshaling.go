package crawler

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

var linePayRegex *regexp.Regexp

func init() {
	//
	reg, err := regexp.Compile(`/^l\d+/`)
	if err != nil {
		fmt.Println("Invalid Regex", err)
		return
	}

	linePayRegex = reg
}

func UnMarshaling(queryFormatData string, uuid string) map[string]interface{} {
	params := strings.Split(queryFormatData, "&")
	resJsonObject := make(map[string]interface{})
	resJsonObject["uuid"] = uuid
	resJsonObject["linepay"] = ``

	for _, paramString := range params {
		keyAndData := strings.Split(paramString, "=")
		key := keyAndData[0]
		value := keyAndData[1]

		if linePayRegex.MatchString(key) {
			linePayValue, ok := resJsonObject["linepay"].(string)
			if !ok {
				log.Fatal("linepay value Error!!!!")
			}
			if len(linePayValue) > 0 {
				linePayValue += ","
			}
			linePayValue += fmt.Sprintf("%s", value)
			resJsonObject["linepay"] = linePayValue
		} else {
			resJsonObject[key] = value
		}
	}
	return resJsonObject
}
