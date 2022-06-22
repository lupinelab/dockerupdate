#!/usr/bin/env python
import os
import stat

# make dockerupdate executeable and move to /usr/local/bin
try:
    os.chmod('dist/dockerupdate', stat.S_IXOTH)
    os.rename('dist/dockerupdate', '/usr/local/bin/dockerupdate')
except PermissionError:
    print("Please run as elevated user")

# move completion script to /etc/bash_completion.d/
try:
    os.rename('dockerupdate_completion', '/etc/bash_completion.d/dockerupdate_completion')
except PermissionError:
    print("Please run as elevated user")