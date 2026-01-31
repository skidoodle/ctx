//go:build darwin

package clipboard

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include <stdlib.h>
#import <Cocoa/Cocoa.h>

int copyFileToPasteboard(char* path) {
    @autoreleasepool {
        NSString *strPath = [NSString stringWithUTF8String:path];
        if (!strPath) return 0;

        NSURL *url = [NSURL fileURLWithPath:strPath];
        if (!url) return 0;

        NSPasteboard *pb = [NSPasteboard generalPasteboard];
        [pb clearContents];

        return [pb writeObjects:@[url]] ? 1 : 0;
    }
}
*/
import "C"
import (
	"errors"
	"path/filepath"
	"unsafe"
)

func CopyFile(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	cPath := C.CString(absPath)
	defer C.free(unsafe.Pointer(cPath))

	if success := C.copyFileToPasteboard(cPath); success == 0 {
		return errors.New("failed to write to macOS pasteboard")
	}
	return nil
}
