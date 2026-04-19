package services

import (
	"regexp"
	"slices"
)

const (
	maxLoginLen        = 64
	maxRoleLen         = 32
	maxIdentifierLen   = 255
	maxFlagLen         = 255
	maxNameLen         = 255
	maxStatusLen       = 32
	maxVariantValueLen = 255
	maxDefaultValueLen = 255
)

var (
	uuidPattern      = regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[1-5][a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`)
	loginPattern     = regexp.MustCompile(`^[A-Za-z0-9._-]{1,64}$`)
	rolePattern      = regexp.MustCompile(`^[A-Za-z0-9_-]{1,32}$`)
	slugPattern      = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_.:-]{0,254}$`)
	statusPattern    = regexp.MustCompile(`^[a-z_]{1,32}$`)
	printablePattern = regexp.MustCompile(`^[[:print:]]+$`)
	statusValues     = []string{"draft", "review", "approved", "active", "pause", "completed", "archive", "rejected"}
)

func isValidUUID(value string) bool {
	return uuidPattern.MatchString(value)
}

func isValidLogin(value string) bool {
	return len(value) <= maxLoginLen && loginPattern.MatchString(value)
}

func isValidRole(value string) bool {
	return len(value) <= maxRoleLen && rolePattern.MatchString(value)
}

func isValidIdentifier(value string) bool {
	return len(value) <= maxIdentifierLen && slugPattern.MatchString(value)
}

func isValidFlag(value string) bool {
	return len(value) <= maxFlagLen && slugPattern.MatchString(value)
}

func isValidStatus(value string) bool {
	return len(value) <= maxStatusLen && statusPattern.MatchString(value) && slices.Contains(statusValues, value)
}

func isValidName(value string) bool {
	return len(value) <= maxNameLen && printablePattern.MatchString(value)
}

func isValidVariantValue(value string) bool {
	return len(value) <= maxVariantValueLen && printablePattern.MatchString(value)
}

func isValidDefaultValue(value string) bool {
	return len(value) <= maxDefaultValueLen && printablePattern.MatchString(value)
}
