package utils

import jsoniter "github.com/json-iterator/go"

// JSON is the globally shared json encoder/decoder
var JSON = jsoniter.ConfigCompatibleWithStandardLibrary
