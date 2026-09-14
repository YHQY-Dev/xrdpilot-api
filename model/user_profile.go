package model

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	AcademicIdentityUndergraduate      = "undergraduate"
	AcademicIdentityMaster             = "master"
	AcademicIdentityPhD                = "phd"
	AcademicIdentityPostdoc            = "postdoc"
	AcademicIdentityLecturer           = "lecturer"
	AcademicIdentityAssociateProfessor = "associate_professor"
	AcademicIdentityProfessor          = "professor"
	AcademicIdentityResearcher         = "researcher"
	AcademicIdentityEngineer           = "engineer"
	AcademicIdentityOther              = "other"

	MaxRealNameLength          = 64
	MaxOrganizationLength      = 128
	MaxSupervisorNameLength    = 64
	MaxResearchDirectionLength = 255
	MaxUsagePurposeLength      = 2000
	maxUsernameAllocateTries   = 1000
)

var validAcademicIdentities = map[string]struct{}{
	AcademicIdentityUndergraduate:      {},
	AcademicIdentityMaster:             {},
	AcademicIdentityPhD:                {},
	AcademicIdentityPostdoc:            {},
	AcademicIdentityLecturer:           {},
	AcademicIdentityAssociateProfessor: {},
	AcademicIdentityProfessor:          {},
	AcademicIdentityResearcher:         {},
	AcademicIdentityEngineer:           {},
	AcademicIdentityOther:              {},
}

var (
	ErrRegistrationEmailRequired    = errors.New("registration email is required")
	ErrRegistrationFieldRequired    = errors.New("registration field is required")
	ErrRegistrationFieldTooLong     = errors.New("registration field is too long")
	ErrRegistrationIdentityInvalid  = errors.New("academic identity is invalid")
	ErrRegistrationUsernameAllocate = errors.New("failed to allocate username")
)

func IsValidAcademicIdentity(code string) bool {
	_, ok := validAcademicIdentities[code]
	return ok
}

func TruncateDisplayName(name string) string {
	return truncateRunes(strings.TrimSpace(name), 20)
}

func UsernameBaseFromEmail(email string) string {
	email = NormalizeEmail(email)
	local, _, _ := strings.Cut(email, "@")
	var b strings.Builder
	prevUnderscore := false
	for i := range len(local) {
		c := local[i]
		isAlphaNum := (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')
		if isAlphaNum {
			b.WriteByte(c)
			prevUnderscore = false
			continue
		}
		if !prevUnderscore {
			b.WriteByte('_')
			prevUnderscore = true
		}
	}
	base := strings.Trim(b.String(), "_")
	if base == "" {
		return "user"
	}
	if len(base) > UserNameMaxLength {
		return base[:UserNameMaxLength]
	}
	return base
}

func UsernameCandidate(base string, attempt int) string {
	if attempt <= 0 {
		if len(base) > UserNameMaxLength {
			return base[:UserNameMaxLength]
		}
		return base
	}
	suffix := strconv.Itoa(attempt + 1)
	maxPrefix := UserNameMaxLength - len(suffix)
	if maxPrefix < 1 {
		maxPrefix = 1
	}
	prefix := base
	if len(prefix) > maxPrefix {
		prefix = prefix[:maxPrefix]
	}
	return prefix + suffix
}

func AllocateUsernameFromEmail(email string) (string, error) {
	base := UsernameBaseFromEmail(email)
	for attempt := range maxUsernameAllocateTries {
		candidate := UsernameCandidate(base, attempt)
		exists, err := CheckUserExistOrDeleted(candidate, "")
		if err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
	}
	return "", ErrRegistrationUsernameAllocate
}

func ValidateRegistrationProfile(user *User) error {
	if user == nil {
		return ErrRegistrationFieldRequired
	}
	user.Email = NormalizeEmail(user.Email)
	user.RealName = strings.TrimSpace(user.RealName)
	user.Organization = strings.TrimSpace(user.Organization)
	user.AcademicIdentity = strings.TrimSpace(user.AcademicIdentity)
	user.SupervisorName = strings.TrimSpace(user.SupervisorName)
	user.ResearchDirection = strings.TrimSpace(user.ResearchDirection)
	user.UsagePurpose = strings.TrimSpace(user.UsagePurpose)

	if user.Email == "" || !looksLikeEmail(user.Email) {
		return ErrRegistrationEmailRequired
	}
	if utf8.RuneCountInString(user.Email) > 50 {
		return ErrRegistrationFieldTooLong
	}
	if err := requireLimitedField(user.RealName, MaxRealNameLength); err != nil {
		return err
	}
	if err := requireLimitedField(user.Organization, MaxOrganizationLength); err != nil {
		return err
	}
	if err := requireLimitedField(user.SupervisorName, MaxSupervisorNameLength); err != nil {
		return err
	}
	if err := requireLimitedField(user.ResearchDirection, MaxResearchDirectionLength); err != nil {
		return err
	}
	if err := requireLimitedField(user.UsagePurpose, MaxUsagePurposeLength); err != nil {
		return err
	}
	if !IsValidAcademicIdentity(user.AcademicIdentity) {
		return ErrRegistrationIdentityInvalid
	}
	return nil
}

func FillPasswordRegistration(user *User) error {
	if user == nil {
		return ErrRegistrationFieldRequired
	}
	user.Username = ""
	if err := ValidateRegistrationProfile(user); err != nil {
		return err
	}
	username, err := AllocateUsernameFromEmail(user.Email)
	if err != nil {
		return err
	}
	user.Username = username
	user.DisplayName = TruncateDisplayName(user.RealName)
	return nil
}

func requireLimitedField(value string, maxRunes int) error {
	if value == "" {
		return ErrRegistrationFieldRequired
	}
	if utf8.RuneCountInString(value) > maxRunes {
		return ErrRegistrationFieldTooLong
	}
	return nil
}

func truncateRunes(value string, maxRunes int) string {
	if maxRunes <= 0 || value == "" {
		return ""
	}
	runes := []rune(value)
	if len(runes) <= maxRunes {
		return value
	}
	return string(runes[:maxRunes])
}

func looksLikeEmail(email string) bool {
	at := strings.IndexByte(email, '@')
	if at <= 0 || at != strings.LastIndexByte(email, '@') {
		return false
	}
	domain := email[at+1:]
	return strings.Contains(domain, ".") && !strings.ContainsAny(email, " \t\r\n")
}
