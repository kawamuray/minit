About
=====
Minimal implementation of /sbin/init which might be enough for Docker container.

Docker is designed to run the target application directly as the pid=1 inside a container.
In most cases it's the valid way to containerize an application, so we don't need ordinary init[1] functionality inside a container.

Though sometimes we still need a init[1] inside a container. For example, when we try to deploy a cluster system which includes functionality to start a daemon on an another node via ssh.
As more concrete example, Hadoop DFS(HDFS) requires ssh running on all slave nodes since it starts datanode daemons by logging into each node via ssh(I know that we can solve this problem by stop using start-dfs.sh and starting each datanode manually if we take more proper way for this case).

Usage
=====
See the result of `minit --help` and [examples/sshd/Dockerfile](https://github.com/kawamuray/minit/blob/master/examples/sshd/Dockerfile).

Capability
==========
- Process supervising       - Monitor orphaned children process and "wait" for them if necessary.
- Syslog                    - Creates /dev/log and accepts connections for syslog logging capability. Output goes to stdout.

Author
======
Yuto Kawamura(kawamuray) <kawamuray.dadada {at} gmail.com>

License
=======
MIT License. Please see LICENSE.
