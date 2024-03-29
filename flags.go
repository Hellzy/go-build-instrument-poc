// Copyright (c) 2016 - 2019 Sqreen. All Rights Reserved.
// Please refer to our terms for more information:
// https://www.sqreen.io/terms.html

package main

import (
	"reflect"
	"strings"
)

type instrumentationToolFlagSet struct {
	Help    bool `ddflag:"-h"`
	Verbose bool `ddflag:"-v"`
	Full    bool `ddflag:"-full"`
}

const structTagKey = "ddflag"

// parseFlags walks through the given arguments and sets the flagSet values
// present in the argument list. Unknown options, not present in the flagSet
// are accepted and skipped. The argument list is not modified.
func parseFlags(flagSet interface{}, args []string) {
	flagSetValueMap := makeFlagSetValueMap(flagSet)

	i := 0
	for i < len(args)-1 {
		_, shift := parseOption(flagSetValueMap, args[i], args[i+1])
		i += shift
	}

	if i < len(args) {
		_, _ = parseOption(flagSetValueMap, args[i], "")
	}
}

func makeFlagSetValueMap(flagSet interface{}) map[string]reflect.Value {
	v := reflect.ValueOf(flagSet).Elem()
	typ := v.Type()
	flagSetValueMap := make(map[string]reflect.Value, v.NumField())
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if tag, ok := field.Tag.Lookup(structTagKey); ok {
			flagSetValueMap[tag] = v.Field(i)
		}
	}
	return flagSetValueMap
}

// parseOption parses the given current argument and following one according to
// the go flags syntax.
func parseOption(flagSetValueMap map[string]reflect.Value, arg, nextArg string) (nonOpt bool, shift int) {
	if arg[0] != '-' {
		// Not an option, return the value and shift by one.
		return true, 1
	}

	// Split the argument by its first `=` character if any, and check the
	// syntax being used.
	kv := strings.SplitN(arg, "=", 2)
	option := kv[0]

	flag, exists := flagSetValueMap[option]

	if len(kv) == 2 {
		// `-opt=val` syntax
		value := kv[1]
		shift = 1
		if exists {
			flag.SetString(value)
		}
	} else if nextArg == "" || len(nextArg) > 1 && nextArg[0] != '-' {
		// `-opt val` syntax
		value := nextArg
		shift = 2
		if exists {
			switch flag.Kind() {
			case reflect.String:
				flag.SetString(value)
			case reflect.Bool:
				flag.SetBool(true)
				shift = 1
			}
		}
	} else {
		// `-opt` syntax (no value)
		shift = 1
		if exists && flag.Kind() == reflect.Bool {
			flag.SetBool(true)
		}
	}

	return
}

func parseFlagsUntilFirstNonOptionArg(flagSet interface{}, args []string) int {
	if len(args) == 0 {
		return -1
	}

	flagSetValueMap := makeFlagSetValueMap(flagSet)

	i := 0
	for i < len(args)-1 {
		nonOpt, shift := parseOption(flagSetValueMap, args[i], args[i+1])
		if nonOpt {
			// First non-option
			return i
		}
		i += shift
	}

	if i < len(args) {
		nonOpt, _ := parseOption(flagSetValueMap, args[i], "")
		if nonOpt {
			// First non-option
			return i
		}
	}

	return -1
}
