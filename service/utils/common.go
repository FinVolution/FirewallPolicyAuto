package utils

import (
	"fmt"
	"net"
)

const PolicyAny = "any"

// slice Remove duplicate elements in slice
func RemoveDuplicateElement[T any](data []T) []T {
	result := make([]T, 0, len(data))
	temp := make(map[any]struct{})

	for _, item := range data {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// CIDRToMask Subnet Mask transfer
func CIDRToMask(cidr int) (string, error) {
	if cidr < 0 || cidr > 32 {
		return "", fmt.Errorf("invalid CIDR prefix length: %d", cidr)
	}

	// get subnet mask 32 bit integer
	// example，cidr = 24，mask = 255.255.255.0
	mask := net.CIDRMask(cidr, 32)

	// 将子网掩码转换为点分十进制格式
	// Convert the subnet mask to dotted decimal format
	ip := net.IP(mask)
	return ip.String(), nil
}

// ContainsAny Assert string list A contains any item in string list B
func ContainsAny(A, B []string) bool {
	for _, itemB := range B {
		for _, itemA := range A {
			if itemB == itemA {
				return true
			}
		}
	}
	return false
}

func AdditionalPolicyItem(policyItemList []string) []string {
	if len(policyItemList) == 0 {
		return []string{PolicyAny}
	}
	return policyItemList
}
