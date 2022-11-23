#!/usr/bin/env python
import os
import shutil
import stat

def install():
    # make dockerupdate executeable and move to /usr/local/bin
    os.chmod('dockerupdate', os.stat('dockerupdate').st_mode | stat.S_IXOTH)
    shutil.copy('dockerupdate', '/usr/local/bin/dockerupdate')
    # move completion script to /etc/bash_completion.d/
    shutil.copy('dockerupdate_completion', '/etc/bash_completion.d/dockerupdate_completion')

if os.getuid() != 0:
    print("Please run as elevated user")
else:
    install()