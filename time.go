package util

import "time"

func ParseTimeForGermany(t string) (time.Time, error) {
	fmt := "2006-01-02 15:04:05"
	return ParseTimeForGermanyFmt(fmt, t)
}
func ParseTimeForGermanyDirName(t string) (time.Time, error) {
	fmt := "2006-01-02_1504"
	return ParseTimeForGermanyFmt(fmt, t)
}
func ParseTimeForGermanyFmt(parseFmt, t string) (time.Time, error) {
	// see time_test for alternative ways to retrieve actual live location
	loc := time.Now().Location()
	tm, err := time.ParseInLocation(parseFmt, "2012-07-31 22:34:45", loc)
	return tm, err
}
