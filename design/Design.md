# Design document to demonstrate various options for transforming data got from text files

## Background

## Achitecture
The idea of the pipelines is layed behind the philosophy of data flows. The data emmitter generates data (in our case files in a proper directories) than consumer process the data and collect them somewhere
So basicly the whole idea looks so:


This idea could be implemented by some different ways
<put glossary>
### One piece of data (POD) with container restart policy
Let's look into the simplest case. We assume that
- service emitting the data and saving it to the files (emitter) located in the same space (VM, Kubernetes Pod, etc). 
- emitter saves the files to virtual file system eg mounted to /tmp folder or any volume in the container - now it doesn't matter
- emitter and consumer are located in the same same space (VM, Kubernetes Pod, etc) and will not be able to work in the separate ones
- The potential bottlenecks here is IO - the disk and the network. But we assume that the speed of the writing file works well and shouldn't be improved and the amount of collected data after processing makes much less than original ones and could be easily handled by the network channel provided us by the cloud
- The durability of the service will be handled by some restart policy for example handled by container restart policy of Kubernetes with exponential backoff