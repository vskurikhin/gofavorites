/*
 * Copyright text:
 * This file was last modified at 2024-07-09 15:02 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * main_test.go
 * $Id$
 */
//!+

package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	flag.Parse()
	exitCode := m.Run()
	// Exit
	os.Exit(exitCode)
}

func TestWithoutArgs(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"main"}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	grpcPort := 65501 + rnd.Intn(30)
	fmt.Fprintf(os.Stderr, "TestWithoutArgs GRPC=127.0.0.1:%d\n", grpcPort)
	t.Setenv("GRPC_ADDRESS", fmt.Sprintf("127.0.0.1:%d", grpcPort))
	httpPort := grpcPort + 1 + rnd.Intn(65535-grpcPort)
	fmt.Fprintf(os.Stderr, "TestWithoutArgs HTTP=127.0.0.1:%d\n", httpPort)
	t.Setenv("HTTP_ADDRESS", fmt.Sprintf("127.0.0.1:%d", httpPort))
	go func() {
		time.Sleep(3 * time.Second)
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()
	main()
}

func BenchmarkTimeSleep(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Sleep(8 * time.Nanosecond)
	}
}

func TestRun(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	run(ctx)
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
