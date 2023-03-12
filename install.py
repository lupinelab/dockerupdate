#!/usr/bin/env python
import os
import shutil
import stat

def install():
    # make dockerupdate executeable and move to /usr/local/bin
    os.chmod('dockerupdate', os.stat('dockerupdate').st_mode | stat.S_IXOTH)
    shutil.copy('dockerupdate', '/usr/local/bin/dockerupdate')

if os.getuid() != 0:
    print("Please run as elevated user")
else:
    install()