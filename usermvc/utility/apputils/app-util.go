package apputils

import (
	"database/sql"
	"encoding/base64"
	"hash/crc32"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

/*
*
 */
func AllSubMatchesWithRegex(submatchRegexp string, input string) []string {

	var retArr []string
	var rExp = regexp.MustCompile(submatchRegexp)
	matches := rExp.FindAllStringSubmatch(input, -1)
	if len(matches) > 0 {
		retArr = make([]string, len(matches))
		for i := range retArr {
			retArr[i] = matches[i][1]
		}
	}
	return retArr
}

/*
*
 */
func AllMatchesWithRegex(matchWithRegex string, input string) []string {

	var rExp = regexp.MustCompile(matchWithRegex)
	return rExp.FindAllString(input, -1)
}

/*
*
 */
func Base64OfImageFromUrl(imgUrl string) (string, error) {

	resp, err := http.Get(imgUrl)
	if err == nil {
		defer resp.Body.Close()
		bytes, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			return base64.StdEncoding.EncodeToString(bytes), nil
		}
	}
	return "", err
}

/*
*
 */
func Base64OfFileContent(filePath string) (string, error) {

	var processingErr error = nil
	var contentStringAsBase64 string = ""
	contentBytes, processingErr := ioutil.ReadFile(filePath)
	if processingErr != nil {
		contentStringAsBase64 = base64.StdEncoding.EncodeToString(contentBytes)
		//replace(/([^\0]{76})/g, "$1\n") + "\n\n";
		//regexp.MustCompile(`([^\0]{76})`).ReplaceAllString(contentStringAsBase64, `$1\n`)
	}
	return contentStringAsBase64, processingErr
}

/*
*
*
 */
func ReplaceKeysWithValues(inText string, keyEncloserStart string,
	keyEncloserEnd string, keyValMap map[string]string) string {

	for key, val := range keyValMap {
		findWhat := keyEncloserStart + key + keyEncloserEnd
		inText = strings.Replace(inText, findWhat, val, -1)
	}
	return inText
}

/*
*
*
 */
func ReplaceKeyWithPrefixedValues(inText string, keyEncloserStart string,
	keyEncloserEnd string, keyValMap map[string]string, prefix string) string {

	for key, val := range keyValMap {
		findWhat := keyEncloserStart + key + keyEncloserEnd
		inText = strings.Replace(inText, findWhat, prefix+val, -1)
	}
	return inText
}

/*
* Does not handle duplicate keys - overwrites
*
 */
func MergeStringMaps(maps ...map[string]string) map[string]string {

	mergedMap := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			mergedMap[k] = v
		}
	}
	return mergedMap
}

/*
*
 */
func StringArrToStringPointersArr(inputArr []string) []*string {

	var retArr []*string
	for i := 0; i < len(inputArr); i++ {
		retArr = append(retArr, &inputArr[i])
	}
	return retArr
}

/*
*
 */
func StringsMapToJsonString(inputMap map[string]string) string {

	var jsonStringBuilder strings.Builder
	jsonStringBuilder.WriteString(`{`)
	for key, element := range inputMap {

		jsonStringBuilder.WriteString(`"`)
		jsonStringBuilder.WriteString(key)
		jsonStringBuilder.WriteString(`":`)
		jsonStringBuilder.WriteString(`"`)
		jsonStringBuilder.WriteString(element)
		jsonStringBuilder.WriteString(`",`)

	}
	tempString := jsonStringBuilder.String()
	return tempString[:len(tempString)-1] + "}"
}

/*
*
 */
func StringArrToCsvString(inputArr []string) string {

	var csvBuilder strings.Builder
	for _, val := range inputArr {
		csvBuilder.WriteString(`"`)
		csvBuilder.WriteString(val)
		csvBuilder.WriteString(`",`)
	}
	retString := csvBuilder.String()
	retString = retString[:len(retString)-1]
	return retString
}

/*
*
 */
func NullColumnValue(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

/*
*
 */
func Crc32OfString(of string) uint32 {
	return crc32.ChecksumIEEE([]byte(of))
}
