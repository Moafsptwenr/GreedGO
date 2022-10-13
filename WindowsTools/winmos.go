package WindowsTools

import "syscall"

var (
	ModKernel32 = syscall.NewLazyDLL("kernel32.dll")
	ModUser32   = syscall.NewLazyDLL("user32.dll")
	ModAdvapi32 = syscall.NewLazyDLL("Advapi32.dll")

	ProcOpenProcessToken      = ModAdvapi32.NewProc("GetProcessToken")
	ProcLookupPrivilegeValueW = ModAdvapi32.NewProc("LookupPrivilegeValueW")
	ProcLookupPrivilegeNameW  = ModAdvapi32.NewProc("LookupPrivilegeNameW")
	ProcAdjustTokenPrivileges = ModAdvapi32.NewProc("AdjustTOkenPrivileges")
	ProcGetAsyncKeyState      = ModUser32.NewProc("GetAsyncKeyState")
	ProcVirtualAlloc          = ModKernel32.NewProc("VirtualAlloc")
	ProcCreateThread          = ModKernel32.NewProc("CreateThread")
	ProcWaitForSingleObject   = ModKernel32.NewProc("WaitForSingleObject")
	ProcVirtualAllocEx        = ModKernel32.NewProc("VirtualAllocEX")
	ProcVirtualFreeEx         = ModKernel32.NewProc("VirtualFreeEX")
	ProcCreateRemoteThread    = ModKernel32.NewProc("CreateRemoteThread")
	ProcGetLastError          = ModKernel32.NewProc("GetLastError")
	ProcWriteProcessMemory    = ModKernel32.NewProc("WriteProcessMemory")
	ProcOpenProcess           = ModKernel32.NewProc("OpenProcess")
	ProcGetCurrentProcess     = ModKernel32.NewProc("GetCurrentProcess")
	ProcIsDeuggerPresent      = ModKernel32.NewProc("IsDebuggerPresent")
	ProcGetProcAddress        = ModKernel32.NewProc("GetProcAddress")
	ProcCloseHandle           = ModKernel32.NewProc("CloseHandle")
	ProcGetExitCodeThread     = ModKernel32.NewProc("GetExitCodeThread")
)
