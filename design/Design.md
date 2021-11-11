# Demonstrating various options for processing large text files

## Achitecture
The idea of the pipelines is layed behind the philosophy of data flows. The data emmitter generates data (in our case files in a proper directories) than consumer process the data and collect them somewhere.

<img src="./img/pipeline.png" width="670">

This idea could be implemented by some different ways
<put glossary>

### One piece of data (POD) with restart policy
Let's look into the simplest case. We assume that
- service emitting the data and saving it to the files (emitter) located in the same space (VM, Kubernetes Pod, etc). 
- emitter saves the files to virtual file system eg mounted to /tmp folder or any volume in the container - now it doesn't matter
- emitter and consumer are located in the same same space (VM, Kubernetes Pod, etc) and will not be able to work in the separate ones
  
<img src="./img/pod-with-restart-policy.png" width="670">

**Potential bottlenecks**

Potential bottleneck here is IO - the disk and the network. But we assume that the speed of the writing file works well and shouldn't be improved and the amount of collected data after processing makes much less than original ones and could be easily handled by the network channel provided us by the cloud provider

**Durability**

The durability of the service will be handled by some restart policy for example by container restart policy of Kubernetes with exponential backoff

**Observability**

Observability of the service will be handled by Prometheus with Alert manager to send alerts to Slack in case of any incident.
Grafana will be used to display statistics about failures that could be used to calculate potential SLO / SLI. 