package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsernameBaseFromEmailUsesLocalPartAndStaysWithinLimit(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  string
	}{
		{
			name:  "plain local part",
			email: "Alice@Example.COM",
			want:  "alice",
		},
		{
			name:  "dots plus and hyphens become underscores",
			email: "jane.doe+lab-1@xrdpilot.com",
			want:  "jane_doe_lab_1",
		},
		{
			name:  "local part longer than 20 is truncated",
			email: "verylonglocalpartnamexyz@xrdpilot.com",
			want:  "verylonglocalpartnam",
		},
		{
			name:  "non ascii local part falls back to user",
			email: "张三@xrdpilot.com",
			want:  "user",
		},
		{
			name:  "empty local part falls back to user",
			email: "@xrdpilot.com",
			want:  "user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UsernameBaseFromEmail(tt.email)
			assert.Equal(t, tt.want, got)
			assert.LessOrEqual(t, len(got), UserNameMaxLength)
		})
	}
}

func TestUsernameCandidateKeepsConflictSuffixWithinLimit(t *testing.T) {
	base := "verylonglocalpartnam"
	got := UsernameCandidate(base, 1)
	assert.Equal(t, "verylonglocalpartna2", got)
	assert.LessOrEqual(t, len(got), UserNameMaxLength)

	got = UsernameCandidate(base, 98)
	assert.Equal(t, "verylonglocalpartn99", got)
	assert.LessOrEqual(t, len(got), UserNameMaxLength)
}

func TestAllocateUsernameFromEmailAddsNumericSuffixWhenTaken(t *testing.T) {
	truncateTables(t)

	require.NoError(t, DB.Create(&User{Username: "alice", Password: "password", AffCode: "aff-alice-1"}).Error)
	require.NoError(t, DB.Create(&User{Username: "alice2", Password: "password", AffCode: "aff-alice-2"}).Error)

	got, err := AllocateUsernameFromEmail("alice@xrdpilot.com")
	require.NoError(t, err)
	assert.Equal(t, "alice3", got)
}

func TestFillPasswordRegistrationIgnoresClientUsernameAndSetsDisplayName(t *testing.T) {
	truncateTables(t)

	user := User{
		Username:          "client_supplied",
		Email:             "  Jane.Doe@Example.COM ",
		Password:          "password12",
		RealName:          "张三丰是一个很长的真实姓名超过二十个字的情况",
		Organization:      "中国科学院",
		AcademicIdentity:  AcademicIdentityPhD,
		SupervisorName:    "李老师",
		ResearchDirection: "计算化学",
		UsagePurpose:      "用于课题计算",
	}

	require.NoError(t, FillPasswordRegistration(&user))

	assert.Equal(t, "jane_doe", user.Username)
	assert.Equal(t, "jane.doe@example.com", user.Email)
	assert.Equal(t, "张三丰是一个很长的真实姓名超过二十个字的情况", user.RealName)
	assert.Equal(t, "张三丰是一个很长的真实姓名超过二十个字的", user.DisplayName)
	assert.LessOrEqual(t, len([]rune(user.DisplayName)), 20)
}

func TestValidateRegistrationProfileRejectsMissingAndInvalidFields(t *testing.T) {
	valid := User{
		Email:             "user@xrdpilot.com",
		RealName:          "张三",
		Organization:      "某大学",
		AcademicIdentity:  AcademicIdentityMaster,
		SupervisorName:    "王老师",
		ResearchDirection: "材料",
		UsagePurpose:      "课程作业",
	}
	require.NoError(t, ValidateRegistrationProfile(&valid))

	missingName := valid
	missingName.RealName = "  "
	require.ErrorIs(t, ValidateRegistrationProfile(&missingName), ErrRegistrationFieldRequired)

	invalidIdentity := valid
	invalidIdentity.AcademicIdentity = "研究生"
	require.ErrorIs(t, ValidateRegistrationProfile(&invalidIdentity), ErrRegistrationIdentityInvalid)

	missingEmail := valid
	missingEmail.Email = ""
	require.ErrorIs(t, ValidateRegistrationProfile(&missingEmail), ErrRegistrationEmailRequired)
}

func TestTruncateDisplayNameCutsByRuneCount(t *testing.T) {
	long := strings.Repeat("研", 25)
	got := TruncateDisplayName(long)
	assert.Equal(t, strings.Repeat("研", 20), got)
	assert.Equal(t, 20, len([]rune(got)))
}
