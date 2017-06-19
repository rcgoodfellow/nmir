# Network Model Intermediate Representation (NMIR)

NMIR is a json based representation for networked system models. It is designed as an intermedieate representation to be generated and consumed by higher level languages. NMIR has the following format.

```json
net:
 id: uuid
 nodes: [node]
 links: [link]
 nets:  [net]

node:
 id: uuid
 endpoints: [endpoint]
 props: {}

link:
 id: uuid
 endpoints: [[uuid],[uuid]]
 props: {}

endpoint:
 id: uuid
 props: {}
```

## net

The net object is a recursive representation of a network. Each network is a collection of nodes, links and subnetworks. Every element in NMIR is identified by a UUID.

## Node

The node object encapsulates a property map and a set of endpoints. Each endpoint can be connected to a link.

## Link

The link object encapsulates a property map and two sets of endpoints. This by having two sets of endpoints instead of simply two endpoints, we retain the ability to model 1:1, 1:\* and \*:\* links. Each endpoint is a uuid reference to an endpoint that is logically owned by a node.

## Endpoint

The endpoint object encapsulates a set of properties. It is the junction between nodes and links.



