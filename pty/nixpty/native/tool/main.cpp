// Note: Due to memory alignment issues, please compile the program separately into
// 32-bit and 64-bit binaries, and then run them to obtain their respective values.

#include <stdio.h>
#include "dragonfly.h"
#include "freebsd.h"
#include "netbsd.h"
#include "openbsd.h"

int main() {
    if (sizeof(int*) == 4) {
        printf("********** 32Bit ARCH **********\n");
    } else {
        printf("********** 64Bit ARCH **********\n");
    }

    dragonfly::output();
    freebsd::output();
    netbsd::output();
    openbsd::output();
    return 0;
}