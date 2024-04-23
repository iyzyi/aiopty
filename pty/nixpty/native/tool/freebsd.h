#include "common.h"

namespace freebsd {
    // copy from FreeBSD/src/sys/sys/filio.h
    struct fiodgname_arg {
        int	    len             // Pay attention to memory alignment
        void	*buf;
    };
    #define	FIODGNAME _IOW('f', 120, struct fiodgname_arg) /* get dev. name */

    void output() {
        printf("[FreeBSD] FIODGNAME: 0x%x\n", FIODGNAME);
    }
}