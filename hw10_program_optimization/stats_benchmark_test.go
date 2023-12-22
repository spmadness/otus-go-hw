//go:build bench
// +build bench

package hw10programoptimization

import (
	"archive/zip"
	"io"
	"testing"
)

func BenchmarkGetUsers(b *testing.B) {
	b.StopTimer()

	data, err := getDataFile()
	if err != nil {
		b.Errorf("data file open error: %v", err)
	}

	defer func(data io.ReadCloser) {
		err = data.Close()
		if err != nil {
			b.Errorf("data file close error: %v", err)
		}
	}(data)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = getUsers(data)
	}
}

func BenchmarkCountDomains(b *testing.B) {
	b.StopTimer()
	data, err := getDataFile()
	if err != nil {
		b.Errorf("data file open error: %v", err)
	}

	defer func(data io.ReadCloser) {
		err = data.Close()
		if err != nil {
			b.Errorf("data file close error: %v", err)
		}
	}(data)

	u, err := getUsers(data)
	if err != nil {
		b.Errorf("get users error: %v", err)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = countDomains(u, "biz")
	}
}

func BenchmarkGetDomainStat(b *testing.B) {
	b.StopTimer()
	data, err := getDataFile()
	if err != nil {
		b.Errorf("data file open error: %v", err)
	}
	defer func(data io.ReadCloser) {
		err = data.Close()
		if err != nil {
			b.Errorf("data file close error: %v", err)
		}
	}(data)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetDomainStat(data, "biz")
	}
}

func getDataFile() (io.ReadCloser, error) {
	r, _ := zip.OpenReader("testdata/users.dat.zip")
	data, err := r.File[0].Open()
	if err != nil {
		return nil, err
	}

	return data, nil
}
