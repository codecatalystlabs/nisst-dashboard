package service

import (
	"testing"
)

func TestMainUploadValidation(t *testing.T) {
	err := validateHeaders("main", []string{
		"SubmissionDate",
		"facility_name",
		"level",
		"region",
		"district",
		"period",
		"meta-instanceID",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
