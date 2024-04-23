#include "common.h"

namespace netbsd {
    // Note: This value is what I defined while somebody may use 1024.
    // See https://github.com/golang/go/issues/66871
    #define PATH_MAX 16

    // copy from NetBSD/src/sys/sys/ttycom.h
    struct ptmget {
        int	    cfd;
        int	    sfd;
        char	cn[PATH_MAX];
        char	sn[PATH_MAX];
    };
    #define TIOCPTSNAME _IOR('t', 72, struct ptmget)	/* ptsname(3) */

    void output() {
        printf("[NetBSD] TIOCPTSNAME: 0x%x\n", TIOCPTSNAME);
    }
}
