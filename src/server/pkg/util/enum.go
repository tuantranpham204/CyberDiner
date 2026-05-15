package util

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/model/entity"
)

var ErrInvalidEnum = errors.New("invalid enum value")

type stringer interface{ ~string }

func ParseEnum[T stringer](raw string, allowed []T) (T, error) {
	s := strings.ToLower(strings.TrimSpace(raw))
	for _, v := range allowed {
		if strings.ToLower(string(v)) == s {
			return v, nil
		}
	}
	var zero T
	return zero, fmt.Errorf("%w: %q not in %s", ErrInvalidEnum, raw, joinEnum(allowed))
}

func IsValidEnum[T stringer](raw string, allowed []T) bool {
	_, err := ParseEnum(raw, allowed)
	return err == nil
}

func joinEnum[T stringer](allowed []T) string {
	parts := make([]string, len(allowed))
	for i, v := range allowed {
		parts[i] = string(v)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

func ParseRole(raw string) (entity.Role, error) {
	return ParseEnum(raw, entity.AllRoles)
}

func ParseGender(raw string) (entity.Gender, error) {
	return ParseEnum(raw, entity.AllGenders)
}

func IsValidRole(raw string) bool   { return IsValidEnum(raw, entity.AllRoles) }
func IsValidGender(raw string) bool { return IsValidEnum(raw, entity.AllGenders) }
