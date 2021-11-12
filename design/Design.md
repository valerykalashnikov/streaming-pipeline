# Demonstrating various options for processing large text files
The idea of the pipelines is layed behind the philosophy of data flows. The data emmitter generates data (in our case files in a proper directories) than consumer process the data and collect them somewhere.

<img src="./img/pipeline.png" width="670">

This idea could be implemented by some different ways
<put glossary>

## One piece of data (POD) with restart policy
Let's look into the simplest case. We assume that
- service emitting the data and saving it to the files (emitter) located in the same space (VM, Kubernetes Pod, etc). 
- emitter saves the files to virtual file system eg mounted to /tmp folder or any volume in the container - now it doesn't matter
- emitter and consumer are located in the same same space (VM, Kubernetes Pod, etc) and will not be able to work in the separate ones

The gathering of the data will be done by periodically scanning the file directory at the virtual filesystem of mounted volume (for example once an hour). The one interesting task is to handle only new files. That means we need to save state after every scanning of the directory. We will accumulate these data in memory and send it to the storage after every scanning.
<img src="./img/pod-with-restart-policy.png" width="670">

The logic of working with every state is:
- scan the whole directory in case of missing data about last scan.
- If the data about last scanning exist, than finding the new files to scan by diffing the list of the files from the directory and list of the files from storage.
- In case of having line number in the state data near the file name, it means that the consumer unexpectedly failed and we need to continue scanning from the next line.

**Potential bottlenecks**

Potential bottleneck here is IO - the disk and the network. But we assume that the speed of the writing file works well and shouldn't be improved and the amount of collected data after processing makes much less than original ones and could be easily handled by the network channel provided us by the cloud provider

**Durability**

The durability of the service will be handled by some restart policy for example by container restart policy of Kubernetes with exponential backoff

**Observability**

Observability of the service will be handled by Prometheus with Alert manager to send alerts to Slack in case of any incident.
Grafana will be used to display statistics about failures that could be used to calculate potential SLO / SLI.

**Trade-ins**
The solution is very simple and could be a good fit for an MVP or in case you don't need a horizontal scaling and very high availability. 
If you need to manage scalability or the network becomes bottlenecks plese look into the next options.

## Worker pool with Ceph storage
In case of necessity to improve scalability we could extend the previous scheme for using multiple workers. Not to handle multiple pieces of disks and keep the simplicity of using the shared one  we could pick some distributed storage like Ceph.

But while introducing scalability we face with the problem of deduplicating data. For example if 2+ workers start processing the same file, the amount of incorrect aggregated data about consumed resource will dramatically increased.  
It could be fixed by handling the state like we explained in the previous option: in case of having line number in the state data near the file name, it means that the consumer unexpectedly failed and we need to continue scanning from the next line.

<img src="./img/ceph-worker-pool.png" width="950">

**Durability**

The durability of the service will be handled by some restart policy for example by container restart policy of Kubernetes with exponential backoff.

Redis could be used as 1 instance because with backup option. Failover of Redises doesn't make sense because in case of failures worker will be able to work on already consumed files and a list of files will not be more than 2Tb ;) In case of really having a necessity to have failover, it's possible to do with Redis. Or to migrate to some distributed K/V like ZooKeeper.

**Potential bottlenecks**
The potential bottlenecks are still Disk IO and potentially could be solved by picking the SSD disks and network in case of using distributed storages.

In case of having network-bloating we could only extend the bottleneck by changing hardware or by using the next option

**Trade-ins**
The solution is more scalabe and fault-taulerant and provides high availability.
From the other hand the trade-in will be potential infrastructure overhead because of necessity to deploy / scale and support worker-services, Ceph cluster and fault-taulerant DBs. 
In case of network is still bottleneck, we could forecast network overpricing in case of using public clouds,or pick the next option

## Worker pool with message queue