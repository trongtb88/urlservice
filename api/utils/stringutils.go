package utils

import (
	"github.com/lucasjones/reggen"
	"log"
	"regexp"
)

func GenerateRandom(regexTemplate string) string {
	g, err := reggen.NewGenerator(regexTemplate)
	if err != nil {
		log.Printf(err.Error())
	}
	return g.Generate(1)
}

func IsMatchRegex(regexTemplate string, value string) bool {
	var isValid = regexp.MustCompile(regexTemplate).MatchString
	return isValid(value)
}


