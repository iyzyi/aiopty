#include "common.h"

namespace dragonfly {
    // copy from DragonFly/src/sys/sys/filio.h
    struct fiodname_args {
    	void	*name;          // Pay attention to memory alignment
    	unsigned int len;
    };
    #define FIODNAME _IOW('f', 120, struct fiodname_args) /* get name of device on that fildesc */

    void output() {
        printf("[DragonFly] FIODNAME: 0x%x\n", FIODNAME);
    }
}