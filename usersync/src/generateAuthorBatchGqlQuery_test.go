package usersync

import (
	"fmt"
	"strings"
	"testing"

	"github.com/OpenQDev/GoGitguru/util/testhelpers"
)

func TestGenerateAuthorBatchGqlQuery(t *testing.T) {
	// ARRANGE - TESTS
	tests := GenerateAuthorBatchGqlQueryTestCases()

	for _, tt := range tests {
		testhelpers.CheckTestSkip(t, testhelpers.Targets(
			"SINGLE_AUTHOR",
		), tt.title)

		t.Run(tt.title, func(t *testing.T) {
			result := generateAuthorBatchGqlQuery(tt.organization, tt.repo, tt.authorList)
			fmt.Println(result)

			sanitizedResult := sanitizeString(result)
			sanitizedExpectedOutput := sanitizeString(tt.expectedOutput)

			if sanitizedResult != sanitizedExpectedOutput {
				t.Errorf("generateAuthorBatchGqlQuery() = %v, want %v", result, tt.expectedOutput)
			}
		})
	}
}

func sanitizeString(str string) string {
	noSpaces := strings.ReplaceAll(str, " ", "")
	noNewLines := strings.ReplaceAll(noSpaces, "\n", "")
	noTabs := strings.ReplaceAll(noNewLines, "\t", "")
	return noTabs
}
