# awsipi
A tool to start/stop AWS OpenShift Cluster's Nodes via command line.

# How to build
**Prerequisite:** Make sure you have a Golang Installed
```
$ git clone https://github.com/anshulvermapatel/awsipi.git
$ cd aswipi
$ go build -o awsipi main.go
$ sudo cp awsipi /usr/local/bin/
$ awsipi help
```
# Do not want to build? 
The built binary is already in the repo. It is just need to copy to any directory in $PATH
```
$ git clone https://github.com/anshulvermapatel/awsipi.git
$ cd aswipi
$ sudo cp awsipi /usr/local/bin/
$ awsipi help
```

# How to Initiate
It is required to be logged in to the AWS OpenShift cluster for which `awsipi` is to be used. This is a one time setup which is to be done, post that it is not required to be logged in to the cluster to start/stop instances of that cluster.
```
- Make sure you are logged in to the AWS OpenShift cluster and execute:
  $ awsipi init <clusterName>
  You can give any <clusterName> for the AWS OpenShift cluster which is logged in.
  Once the cluster is initialized, a directory with name `.awsipi` will be created in the current user's home directory. The directory will contain a file of the name of the <clusterName> with which it was initialized. That file will contain the instance IDs and zone of the AWS OpenShift cluster's instances.
  
  Once the above is done. The instances can be simply Started/Stopped by executing the following command.
  To start the instances -
  $ awsipi start <clusterName>
  $ awsipi stop <clusterName>
```
