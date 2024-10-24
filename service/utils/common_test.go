package utils

import (
	"fmt"
	"testing"
)

// TestRemoveDuplicateElement 测试 RemoveDuplicateElement 函数
func TestRemoveDuplicateElement(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "Remove duplicates from slice",
			input:    []int{1, 2, 2, 3, 4, 4, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "Empty slice",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "Slice with no duplicates",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := RemoveDuplicateElement(test.input)
			if len(actual) != len(test.expected) {
				t.Errorf("Expected %v, got %v", test.expected, actual)
				return
			}
			for i := range actual {
				if actual[i] != test.expected[i] {
					t.Errorf("Expected %v, got %v", test.expected, actual)
					return
				}
			}
		})
	}
}

// TestCIDRToMask 测试 CIDRToMask 函数
func TestCIDRToMask(t *testing.T) {
	tests := []struct {
		cidr     int
		expected string
	}{
		{0, "0.0.0.0"},
		{8, "255.0.0.0"},
		{16, "255.255.0.0"},
		{24, "255.255.255.0"},
		{32, "255.255.255.255"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("CIDR%d", test.cidr), func(t *testing.T) {
			mask, err := CIDRToMask(test.cidr)
			if err != nil {
				t.Fatalf("CIDRToMask(%d) error = %v", test.cidr, err)
			}
			if mask != test.expected {
				t.Errorf("CIDRToMask(%d) = %s, want %s", test.cidr, mask, test.expected)
			}
		})
	}
}

// TestContainsAny 测试 ContainsAny 函数
func TestContainsAny(t *testing.T) {
	tests := []struct {
		A        []string
		B        []string
		expected bool
	}{
		{[]string{"a", "b", "c"}, []string{"b"}, true},
		{[]string{"a", "b", "c"}, []string{"d"}, false},
		{[]string{"a", "b", "c"}, []string{}, false},
		{[]string{}, []string{"a"}, false},
	}

	for _, test := range tests {
		t.Run("ContainsAny", func(t *testing.T) {
			actual := ContainsAny(test.A, test.B)
			if actual != test.expected {
				t.Errorf("ContainsAny(%v, %v) = %v, want %v", test.A, test.B, actual, test.expected)
			}
		})
	}
}

// TestAdditionalPolicyItem 测试 AdditionalPolicyItem 函数
func TestAdditionalPolicyItem(t *testing.T) {
	tests := []struct {
		policyItemList []string
		expected       []string
	}{
		{[]string{}, []string{PolicyAny}},
		{[]string{PolicyAny}, []string{PolicyAny}},
		{[]string{"policy1", "policy2"}, []string{"policy1", "policy2"}},
	}

	for _, test := range tests {
		t.Run("AdditionalPolicyItem", func(t *testing.T) {
			actual := AdditionalPolicyItem(test.policyItemList)
			if len(actual) != len(test.expected) {
				t.Errorf("Expected %v, got %v", test.expected, actual)
				return
			}
			for i := range actual {
				if actual[i] != test.expected[i] {
					t.Errorf("Expected %v, got %v", test.expected, actual)
					return
				}
			}
		})
	}
}
