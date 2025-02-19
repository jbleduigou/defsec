package test

import (
	"testing"

	"github.com/aquasecurity/defsec/test/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_IsPresentCheckOnBlock(t *testing.T) {
	var tests = []struct {
		name              string
		source            string
		expectedAttribute string
	}{
		{
			name: "expected attribute is present",
			source: `
resource "aws_s3_bucket" "my-bucket" {
 	bucket_name = "bucketName"
}`,
			expectedAttribute: "bucket_name",
		},
		{
			name: "expected acl attribute is present",
			source: `
resource "aws_s3_bucket" "my-bucket" {
 	bucket_name = "bucketName"
	acl = "public-read"
}`,
			expectedAttribute: "acl",
		},
		{
			name: "expected acl attribute is present",
			source: `
resource "aws_s3_bucket" "my-bucket" {
 	bucket_name = "bucketName"
	acl = "public-read"
	logging {
		target_bucket = aws_s3_bucket.log_bucket.id
		target_prefix = "log/"
	}
}`,
			expectedAttribute: "logging",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			modules := testutil.CreateModulesFromSource(test.source, ".tf", t)
			for _, module := range modules {
				for _, block := range module.GetBlocks() {
					assert.Equal(t, block.HasChild(test.expectedAttribute), true)
					assert.Equal(t, !block.HasChild(test.expectedAttribute), false)
				}
			}
		})
	}
}

func Test_IsNotPresentCheckOnBlock(t *testing.T) {
	var tests = []struct {
		name              string
		source            string
		expectedAttribute string
	}{
		{
			name: "expected attribute is not present",
			source: `
resource "aws_s3_bucket" "my-bucket" {
 	bucket_name = "bucketName"
	
}`,
			expectedAttribute: "acl",
		},
		{
			name: "expected acl attribute is not present",
			source: `
resource "aws_s3_bucket" "my-bucket" {
 	bucket_name = "bucketName"
	acl = "public-read"
	
}`,
			expectedAttribute: "logging",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			modules := testutil.CreateModulesFromSource(test.source, ".tf", t)
			for _, module := range modules {
				for _, block := range module.GetBlocks() {
					assert.Equal(t, block.HasChild(test.expectedAttribute), false)
					assert.Equal(t, !block.HasChild(test.expectedAttribute), true)
				}
			}
		})
	}
}

func Test_MissingChildNotFoundOnBlock(t *testing.T) {
	var tests = []struct {
		name              string
		source            string
		expectedAttribute string
	}{
		{
			name: "expected attribute is not present",
			source: `
resource "aws_s3_bucket" "my-bucket" {
 	bucket_name = "bucketName"
	
}`,
			expectedAttribute: "acl",
		},
		{
			name: "expected acl attribute is not present",
			source: `
resource "aws_s3_bucket" "my-bucket" {
 	bucket_name = "bucketName"
	acl = "public-read"
	
}`,
			expectedAttribute: "logging",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			modules := testutil.CreateModulesFromSource(test.source, ".tf", t)
			for _, module := range modules {
				for _, block := range module.GetBlocks() {
					assert.Equal(t, block.MissingChild(test.expectedAttribute), true)
					assert.Equal(t, !block.HasChild(test.expectedAttribute), true)
				}
			}
		})
	}
}
