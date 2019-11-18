package main

import (
    "regexp"
)

func oodle(s string) string {
    re_lower := regexp.MustCompile(`[aeiou]`)
    re_upper := regexp.MustCompile(`[AEIOU]`)

    nstr := re_lower.ReplaceAllString(s, "oodle")
    nstr = re_upper.ReplaceAllString(nstr, "OODLE")

    return nstr
}
