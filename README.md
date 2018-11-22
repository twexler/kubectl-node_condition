kubectl node-condition
======================

`kubectl node-nodition` is a plugin for [`kubectl`](https://kubernetes.io/docs/reference/kubectl/overview/) that simply outputs all Conditions linked to a Node.

Install
-------

``` shell
go get github.com/twexler/kube-node_condition
```

Example
-------

``` shell
> kubectl node-condition
docker-for-desktop
==================

+----------------------------+--------+--------------------------------+-------------------------------+
|           REASON           | STATUS |            MESSAGE             |     LAST TRANSITION TIME      |
+----------------------------+--------+--------------------------------+-------------------------------+
| KubeletHasSufficientDisk   | False  | kubelet has sufficient disk    | 2018-11-22 12:17:49 -0500 EST |
|                            |        | space available                |                               |
+----------------------------+--------+--------------------------------+-------------------------------+
| KubeletHasSufficientMemory | False  | kubelet has sufficient memory  | 2018-11-22 12:17:49 -0500 EST |
|                            |        | available                      |                               |
+----------------------------+--------+--------------------------------+-------------------------------+
| KubeletHasNoDiskPressure   | False  | kubelet has no disk pressure   | 2018-11-22 12:17:49 -0500 EST |
+----------------------------+--------+--------------------------------+-------------------------------+
| KubeletHasSufficientPID    | False  | kubelet has sufficient PID     | 2018-11-22 12:17:49 -0500 EST |
|                            |        | available                      |                               |
+----------------------------+--------+--------------------------------+-------------------------------+
| KubeletReady               | True   | kubelet is posting ready       | 2018-11-22 12:17:49 -0500 EST |
|                            |        | status                         |                               |
+----------------------------+--------+--------------------------------+-------------------------------+

```

Why
---

I found it difficult to ascertain the condition of nodes on clusters I operate en masse.