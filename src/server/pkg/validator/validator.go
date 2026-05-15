package validator

import (
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/tuantranpham204/CyberDiner.git/src/server/pkg/util"
)

var (
	usernameRe = regexp.MustCompile(`^[a-zA-Z0-9._-]{3,30}$`)
	nameRe     = regexp.MustCompile(`^[\p{L}][\p{L}\s'-]{0,49}$`)
	// Vietnamese mobile/landline: either "0" or "+84" followed by 9 digits.
	// Spaces, dashes, dots, and parentheses in the input are stripped before
	// this pattern is applied (see phone()).
	vnPhoneRe = regexp.MustCompile(`^(?:\+84|0)\d{9}$`)
)

var phoneStripper = strings.NewReplacer(
	" ", "",
	"-", "",
	"(", "",
	")", "",
	".", "",
)

func Register(v *validator.Validate) error {
	if err := v.RegisterValidation("username", username); err != nil {
		return err
	}
	if err := v.RegisterValidation("strongpwd", strongPassword); err != nil {
		return err
	}
	if err := v.RegisterValidation("gender", gender); err != nil {
		return err
	}
	if err := v.RegisterValidation("role", role); err != nil {
		return err
	}
	if err := v.RegisterValidation("dob", dob); err != nil {
		return err
	}
	if err := v.RegisterValidation("personname", personName); err != nil {
		return err
	}
	if err := v.RegisterValidation("phone", phone); err != nil {
		return err
	}
	return nil
}

func username(fl validator.FieldLevel) bool {
	return usernameRe.MatchString(fl.Field().String())
}

func personName(fl validator.FieldLevel) bool {
	return nameRe.MatchString(fl.Field().String())
}

func strongPassword(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	if len(s) < 8 || len(s) > 72 {
		return false
	}
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, r := range s {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasDigit && hasSpecial
}

func gender(fl validator.FieldLevel) bool {
	return util.IsValidGender(fl.Field().String())
}

func role(fl validator.FieldLevel) bool {
	return util.IsValidRole(fl.Field().String())
}

func dob(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return false
	}
	now := time.Now()
	if t.After(now) {
		return false
	}
	age := now.Year() - t.Year()
	if now.YearDay() < t.YearDay() {
		age--
	}
	return age >= 13 && age <= 120
}

// phone validates a Vietnamese phone number. Accepts common human-friendly
// formats by stripping spaces, dashes, dots, and parentheses before matching,
// so "0912 345 678", "0912-345-678", and "+84 912 345 678" are all valid.
func phone(fl validator.FieldLevel) bool {
	raw := fl.Field().String()
	normalized := phoneStripper.Replace(raw)
	return vnPhoneRe.MatchString(normalized)
}
