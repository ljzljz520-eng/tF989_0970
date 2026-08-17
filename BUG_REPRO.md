# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	followuplabel/cmd/app	[no test files]
ok  	followuplabel/internal/controller	0.002s
?   	followuplabel/internal/model	[no test files]
?   	followuplabel/internal/repository	[no test files]
--- FAIL: TestUninstallWithRetentionRemovesMenuAndPreservesTags (0.00s)
panic: runtime error: invalid memory address or nil pointer dereference [recovered]
	panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x125ee8]

goroutine 8 [running]:
testing.tRunner.func1.2({0x159b80, 0x2a5950})
	/usr/local/go/src/testing/testing.go:1734 +0x1ac
testing.tRunner.func1()
	/usr/local/go/src/testing/testing.go:1737 +0x334
panic({0x159b80?, 0x2a5950?})
	/usr/local/go/src/runtime/panic.go:792 +0x124
followuplabel/internal/service.(*Plugin).Uninstall(0x400008de10, {0x1bd310, 0x2d6060}, 0x1)
	/app/internal/service/plugin.go:46 +0x58
followuplabel/internal/service_test.TestUninstallWithRetentionRemovesMenuAndPreservesTags(0x40000ca000)
	/app/internal/service/plugin_test.go:63 +0x484
testing.tRunner(0x40000ca000, 0x18f490)
	/usr/local/go/src/testing/testing.go:1792 +0xe4
created by testing.(*T).Run in goroutine 1
	/usr/local/go/src/testing/testing.go:1851 +0x374
FAIL	followuplabel/internal/service	0.005s
?   	followuplabel/internal/validation	[no test files]
?   	followuplabel/web	[no test files]
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/app): exit `0`
- Frontend build (web): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/app): exit `0`
- Frontend build (web): exit `0`
