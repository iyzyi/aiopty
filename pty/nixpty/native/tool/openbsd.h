#include "common.h"

namespace openbsd {
    // copy from OpenBSD/src/sys/sys/tty.h
    struct ptmget {
        int	    cfd;
        int	    sfd;
        char	cn[16];
        char	sn[16];
    };
    #define PTMGET _IOR('t', 1, struct ptmget)

    void output() {
        printf("[OpenBSD] PTMGET: 0x%x\n", PTMGET);
    }
}
