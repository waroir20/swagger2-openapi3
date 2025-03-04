package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
)

func ConvertSwaggerSpec2OpenAPI3dot0(fileIn, fileOut string) {
	file, err := os.ReadFile(fileIn)
	if err != nil {
		fmt.Printf("Error read from file: '%s'\n", fileIn)
		return
	}
	stringified := string(file)
	stringified = strings.Replace(stringified, "#/definitions/", "#/components/schemas/", -1)
	var rawJSONMap map[string]interface{}
	err = json.Unmarshal([]byte(stringified), &rawJSONMap)
	if err != nil {
		fmt.Printf("Count not read schema for %s: %s", fileIn, err)
		return
	}
	rawJSONMap["openapi"] = "3.0.0"
	delete(rawJSONMap, "swagger")
	delete(rawJSONMap, "schemes")
	delete(rawJSONMap, "servers")
	components := map[string]interface{}{}

	components["schemas"] = rawJSONMap["definitions"]
	bearer, ok := getSecuritySchemeBearer(rawJSONMap)
	if ok {
		components["securitySchemes"] = map[string]interface{}{
			"Bearer": bearer,
		}
	}
	rawJSONMap["components"] = components
	delete(rawJSONMap, "definitions")
	delete(rawJSONMap, "securityDefinitions")
	delete(rawJSONMap, "host")
	delete(rawJSONMap, "basePath")
	// rawJSONMap["servers"] = rawJSONMapRef["servers"]

	if paths, ok := rawJSONMap["paths"]; ok {
		jsonMap := paths.(map[string]interface{})
		replaceRequestDetails(jsonMap)
		replaceResponseDetails(jsonMap)
	}

	schema, err := json.MarshalIndent(rawJSONMap, "", "  ")
	if err != nil {
		fmt.Printf("Count not re-generate schema %s", err)
		return
	}
	err = os.WriteFile(fileOut, schema, 0644)
	if err != nil {
		fmt.Printf("Count not write schema to file for %s: %s", fileOut, err)
		return
	}
}

func getSecuritySchemeBearer(rawJSONMap map[string]interface{}) (map[string]interface{}, bool) {
	secDefs, ok := rawJSONMap["securityDefinitions"]
	if !ok {
		return rawJSONMap, false
	}
	secDef, ok := secDefs.(map[string]interface{})["Bearer"]
	var newSecDef map[string]interface{}
	if ok {
		newSecDef = make(map[string]interface{})
		secDefMap := secDef.(map[string]interface{})
		newSecDef["description"] = secDefMap["description"]
		newSecDef["bearerFormat"] = "Bearer"
		newSecDef["scheme"] = "bearer"
		newSecDef["type"] = "http"
	}
	return newSecDef, newSecDef != nil
}

func replaceResponseDetails(rawJSONMap map[string]interface{}) {
	foundResponsesLevel := false
	responses := interface{}(nil)
	for key, value := range rawJSONMap {
		if key == "responses" {
			responses = value
			foundResponsesLevel = true
		}
	}
	if !foundResponsesLevel {
		for _, value := range rawJSONMap {
			replaceResponseDetails(value.(map[string]interface{}))
		}
		return
	}
	contentType := interface{}("application/json")
	produces, ok := rawJSONMap["produces"]
	if ok {
		contentTypes := produces.([]interface{})
		if len(contentTypes) > 0 {
			contentType = contentTypes[0]
		}
	}
	replaceResponseSchema(responses.(map[string]interface{}), fmt.Sprintf("%v", contentType))
	delete(rawJSONMap, "produces")
}

func replaceResponseSchema(rawJSONMap map[string]interface{}, contentType string) {
	foundSchemaLevel := false
	schema := interface{}(nil)
	for key, value := range rawJSONMap {
		if key == "schema" {
			schema = value
			foundSchemaLevel = true
		}
	}
	if !foundSchemaLevel {
		kind := reflect.ValueOf(map[string]interface{}{}).Kind()
		for _, value := range rawJSONMap {
			if reflect.ValueOf(value).Kind() == kind {
				replaceResponseSchema(value.(map[string]interface{}), contentType)
			}
		}
		return
	}

	newDef := map[string]map[string]interface{}{
		contentType: {
			"schema": schema,
		},
	}
	delete(rawJSONMap, "schema")
	rawJSONMap["content"] = newDef
}

func replaceRequestDetails(rawJSONMap map[string]interface{}) {
	foundConsumesLevel := false
	consumes := interface{}(nil)
	for key, value := range rawJSONMap {
		if key == "consumes" {
			consumes = value
			foundConsumesLevel = true
		}
	}
	if !foundConsumesLevel {
		kind := reflect.ValueOf(map[string]interface{}{}).Kind()
		for _, value := range rawJSONMap {
			if reflect.ValueOf(value).Kind() == kind {
				replaceRequestDetails(value.(map[string]interface{}))
			}
		}
		return
	}
	contentType := interface{}("application/json")
	contentTypes := consumes.([]interface{})
	if len(contentTypes) > 0 {
		contentType = contentTypes[0]
	}
	delete(rawJSONMap, "consumes")

	parameters, ok := rawJSONMap["parameters"]
	if ok {
		params := parameters.([]interface{})
		newParams := make([]map[string]interface{}, 0)
		requestBody := interface{}(nil)
		for _, oldParam := range params {
			old := oldParam.(map[string]interface{})
			sch, ok := old["schema"]
			in, okIn := old["in"]
			if ok && okIn && fmt.Sprintf("%v", in) == "body" {
				requestBody = sch
				continue
			}
			newParam := map[string]interface{}{}
			if o, ok := old["description"]; ok {
				newParam["description"] = o
			}
			if o, ok := old["name"]; ok {
				newParam["name"] = o
			}
			if o, ok := old["required"]; ok {
				newParam["required"] = o
			}
			if o, ok := old["in"]; ok {
				newParam["in"] = o
			}
			mappy := map[string]interface{}{}
			if typ, ok := old["type"]; ok {
				mappy["type"] = typ
			}
			if enum, ok := old["enum"]; ok {
				mappy["enum"] = enum
			}
			if len(mappy) > 0 {
				newParam["schema"] = mappy
			}
			newParams = append(newParams, newParam)
		}

		delete(rawJSONMap, "parameters")
		rawJSONMap["parameters"] = newParams
		if requestBody != nil {
			rawJSONMap["requestBody"] = map[string]map[string]map[string]interface{}{
				"content": {
					fmt.Sprintf("%v", contentType): {
						"schema": requestBody,
					},
				},
			}
		}
	}
}
