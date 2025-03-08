# GONE-CLI

GONE-CLI is a go application to interact with [GONE](https://github.com/David-Antunes/gone) and [GONE-Agent](https://github.com/David-Antunes/gone-agent).

```
CLI tool to interact with GONE

Usage:
  gone-cli [command]

Available Commands:
  agent       Manage interaction with gone-agent
  bridge      Adds a bridge to emulation
  connect     Connects network components together.
  disconnect  Disconnects components from the emulation
  disrupt     Disrupts network components in the network emulation
  forget      Cleans routing rules for a specific router
  help        Help about any command
  inspect     Shows information about the emulation
  intercept   Intercepts traffic from a given Link
  network     Controls whether a particular bridge or router executes
  node        Adds node to the network emulation
  pause       Pauses a node
  propagate   A brief description of your command
  remove      Removes a component from the network emulation
  router      Adds new router to the network emulation.
  sniff       Sniffs traffic from Link
  unpause     Unpauses a node

Flags:
      --curl      Prints curl request
      --dry-run   Does not make request. Useful for generating json request or curl request.
  -h, --help      help for gone-cli
      --json      Prints json request
      --quiet     Supresses the sucessfull Response

Use "gone-cli [command] --help" for more information about a command.
```

## Environment Variables

GONE-CLI by default, sends its requests to localhost. However, the user can change the endpoint through environment variables.

### Example

```bash
GONE_URL=http://localhost:3000 AGENT_URL=http://localhost:3300 gone-cli 
```

## Commands

At its core, this tool only converts the desired action into the correct HTTP request. In the following sections, we detail the command used, the JSON body, and the equivalent curl request and the response.

Although we present the the json requests, curl request, and json response, there are commands that utilize other endpoints besides the ones shown.

To get the information about the request and curl request, the user can use the flags "--json" and "--curl" flags for GONE-CLI to output the created json body and equivalent curl request.

### Node

Adds a new application to the network emulator.

```bash
gone-cli node -- docker run -d --network gone_net --name ubuntu2 ubuntu sleep 100000
```

```json
{
  "dockerCmd": [
    "docker",
    "run",
    "-d",
    "--network",
    "gone_net",
    "--name",
    "ubuntu2",
    "ubuntu",
    "sleep",
    "10000"
  ],
  "machineId": ""
}
```

To add a new application to the emulation, the user must provide the docker command after the "--" and must contain the "-d" and "--network gone_net" fields in the command. The command fails otherwise.

CURL:

```bash
curl -X 'POST' -d '{"dockerCmd":["docker","run","-d","--network","gone_net","--name","ubuntu2","ubuntu","sleep","10000"],"machineId":""}' -H 'Content-Type: application/json' 'http://localhost:3000/addNode'
```

RESPONSE:

```json
{
  "name": "ubuntu2",
  "mac": "02:42:0a:01:00:b0",
  "ip": "10.1.0.176",
  "machineId": "primary",
  "err": {
    "err_code": 0,
    "err_msg": ""
  }
}
```

### Bridge

To add a new bridge, the user must supply a new id.

```bash
gone-cli bridge bridge1
```

```json
{
  "name": "bridge1",
  "machineId": ""
}
```

CURL:

```bash
curl -X 'POST' -d '{"name":"bridge1","machineId":""}' -H 'Content-Type: application/json' 'http://localhost:3000/addBridge'
```

RESPONSE:

```json
{
  "name": "bridge1",
  "machineId": "primary",
  "err": {
    "err_code": 0,
    "err_msg": ""
  }
}
```

### Router

Adds a new router to the emulation

```bash
gone-cli router router1
```

```json
{
  "name": "router1",
  "machineId": ""
}
```

CURL:

```bash
curl -X 'POST' -d '{"name":"router1","machineId":""}' -H 'Content-Type: application/json' 'http://localhost:3000/addRouter'
```

RESPONSE:

```json
{
  "name": "router1",
  "machineId": "primary",
  "err": {
    "err_code": 0,
    "err_msg": ""
  }
}
```

### Connect

This command allows connecting the components in the emulation by defining some network properties.

```bash
gone-cli --json --curl connect -l 10.0 -w 10M -d 0.01 -j 1.0 -c 100 -b bridge1 router1
```

```json
{
  "bridge": "bridge1",
  "router": "router1",
  "latency": 10,
  "jitter": 1,
  "dropRate": 0.01,
  "bandwidth": 10000000,
  "weight": 0
}
```

The configurable network properties in the emulation are: latency (-l), bandwidth (-w), packet loss (-d), jitter (-j), link weight (-c).

Depending on the connecting components, the user must choose the correct connection.

1. "-n" requires a valid node id and bridge id
2. "-b" requires a valid bridge id and router id
3. "-r" requires two valid router ids

CURL:

```bash
curl -X 'POST' -d '{"bridge":"bridge1","router":"router1","latency":10,"jitter":1,"dropRate":0.01,"bandwidth":10000000,"weight":0}' -H 'Content-Type: application/json' 'http://localhost:3000/connectBridgeToRouter'
```

RESPONSE:

```json
{
  "bridge": "bridge1",
  "err": {
    "err_code": 0,
    "err_msg": ""
  },
  "router": "router1"
}
```

### Disconnect

The user can also disconnect components.

```bash
gone-cli disconnect -b bridge
```

```json
{
  "name": "bridge1"
}
```
Depending on link to disconnect, the user must provide the correct flag.

1. "-n" requires a valid node id and disconnects the node from its bridge
2. "-b" requires a valid bridge id and disconnects the bridge from its router
3. "-r" requires two valid router ids

CURL:

```bash
curl -X 'POST' -d '{"name":"bridge1"}' -H 'Content-Type: application/json' 'http://localhost:3000/disconnectBridge'
```

RESPONSE:

```json
{
  "err": {
    "err_code": 0,
    "err_msg": ""
  },
  "name": "bridge1"
}
```

### Inspect

The user can inspect any component and obtain information about their connections.

```bash
gone-cli inspect -n ubuntu2
```

```json
{
  "name": "ubuntu2"
}
```

CURL:
```bash
curl -X 'POST' -d '{"name":"ubuntu2"}' -H 'Content-Type: application/json' 'http://localhost:3000/inspectNode'
```

RESPONSE:

```json
{
  "node": {
    "Id": "ubuntu2",
    "Mac": "02:42:0a:01:00:b0",
    "MachineId": "primary",
    "Bridge": "bridge1",
    "Link": {
      "To": "ubuntu2",
      "From": "bridge1",
      "LinkProps": {
        "Latency": 0,
        "Bandwidth": 1250000,
        "Jitter": 0,
        "DropRate": 0,
        "Weight": 0
      }
    }
  },
  "err": {
    "err_code": 0,
    "err_msg": ""
  }
}
```

### Remove

The user can also remove a component from the emulation.

```bash
gone-cli remove -r router2
```

```json
{
  "name": "router2"
}
```

CURL:

```bash
curl -X 'POST' -d '{"name":"router2"}' -H 'Content-Type: application/json' 'http://localhost:3000/removeRouter'
```

RESPONSE:

```json
{
  "err": {
    "err_code": 0,
    "err_msg": ""
  },
  "name": "router2"
}
```

### Pause

The user can pause a node or all the nodes.

```bash
gone-cli pause ubuntu1
```

```json
{
  "id": "ubuntu1",
  "all": false
}
```

CURL:

```bash
curl -X 'POST' -d '{"id":"ubuntu1","all":false}' -H 'Content-Type: application/json' 'http://localhost:3000/pause'
```

RESPONSE:

```json
{
  "id": "ubuntu1",
  "all": false,
  "error": {
    "err_code": 0,
    "err_msg": ""
  }
}
```

### Unpause

The user can unpause a node or all the nodes.

```bash
gone-cli unpause ubuntu1
```

```json
{
  "id": "ubuntu1",
  "all": false
}
```

CURL:

```bash
curl -X 'POST' -d '{"id":"ubuntu1","all":false}' -H 'Content-Type: application/json' 'http://localhost:3000/unpause'
```

RESPONSE:

```json
{
  "id": "ubuntu1",
  "all": false,
  "error": {
    "err_code": 0,
    "err_msg": ""
  }
}
```

### Propagate

To optimize the emulated network, the user can trigger a propagation of the routing rules of a specific router.

```bash
gone-cli propagate router1
```

```json
{
  "name": "router1"
}
```

CURL:

```bash
curl -X 'POST' -d '{"name":"router1"}' -H 'Content-Type: application/json' 'http://localhost:3000/propagate'
```

RESPONSE:

```json
{
  "name": "router1",
  "errors": {
    "err_code": 0,
    "err_msg": ""
  }
}
```

### Forget

The user can also make a router forget its routing rules.

```bash
gone-cli forget router1
```

```json
{
  "name": "router1"
}
```

CURL:

```bash
curl -X 'POST' -d '{"name":"router1"}' -H 'Content-Type: application/json' 'http://localhost:3000/forget'
```

RESPONSE:

```json
{
  "name": "router1",
  "errors": {
    "err_code": 0,
    "err_msg": ""
  }
}
```

### Disrupt

The user can temporarily disable a link.

```bash
gone-cli disrupt -b bridge1
```


```json
{
  "bridge": "bridge1"
}
```
Depending on link to disrupt, the user must provide the correct flag.

1. "-n" requires a valid node id and disrupts the link between node and bridge
2. "-b" requires a valid bridge id and disrupts the link between the bridge and router
3. "-r" requires two valid router ids and disrupts their connection

CURL: 

```bash
curl -X 'POST' -d '{"bridge":"bridge1"}' -H 'Content-Type: application/json' 'http://localhost:3000/disruptBridge'
```

RESPONSE:

```json
{
  "error": {
    "err_code": 0,
    "err_msg": ""
  },
  "router1": "router1",
  "router2": "router2"
}
```

### Network

The user can also temporarily disable a bridge or a router.

```bash
gone-cli network -b -s bridge1
```

```json
{
  "bridge": "bridge1"
}
```

Not providing the "-s" flag starts the component.

CURL:

```bash
curl -X 'POST' -d '{"bridge":"bridge1"}' -H 'Content-Type: application/json' 'http://localhost:3000/stopBridge'
```

RESPONSE:

```json
{
  "bridge": "bridge1",
  "error": {
    "err_code": 0,
    "err_msg": ""
  }
}
```

### Sniff

The user can receive a copy of the network traffic of a given link.

```bash
gone-cli sniff -i link1 -b bridge1
```

```json
{
  "bridge": "bridge1",
  "id": "link1"
}
```
Depending on link to sniff, the user must provide the correct flag.

1. "-n" requires a valid node id and sniffs the link between node and bridge
2. "-b" requires a valid bridge id and sniffs the link between the bridge and router
3. "-r" requires two valid router ids and sniffs their connection
4. "-s" requires an id to stop the sniffing

If the user does not provide the "-i" flag to create a custom id for the socket, a unique id is generated instead.

CURL:

```bash
curl -X 'POST' -d '{"bridge":"bridge1","id":"link1"}' -H 'Content-Type: application/json' 'http://localhost:3000/sniffBridge'
```

```json
{
  "bridge": "bridge1",
  "errors": {
    "err_code": 0,
    "err_msg": ""
  },
  "id": "link1",
  "machineId": "primary",
  "path": "/tmp/link1.sniff"
}
```

There is an example over on [GONE-Sniffer](https://github.com/David-Antunes/gone-sniffer) leveraging this type of operation.

### Intercept

The user instead of receiving a copy, can completely intercept the network traffic.

```bash
gone-cli intercept -i link2 -r router1 router2
```

```json
{
  "router1": "router1",
  "router2": "router2",
  "id": "link2",
  "direction": true
}
```
Depending on link to intercept, the user must provide the correct flag.

1. "-n" requires a valid node id and intercepts the link between node and bridge
2. "-b" requires a valid bridge id and intercepts the link between the bridge and router
3. "-r" requires two valid router ids and intercepts their connection
4. "-s" requires an id to stop the interception 

If the user does not provide the "-i" flag to create a custom id for the socket, a unique id is generated instead.

CURL:

```bash
curl -X 'POST' -d '{"router1":"router1","router2":"router2","id":"link2","direction":true}' -H 'Content-Type: application/json' 'http://localhost:3000/interceptRouter'
```

RESPONSE:

```json
{
  "error": {
    "err_code": 0,
    "err_msg": ""
  },
  "id": "link2",
  "machineId": "primary",
  "path": "/tmp/link2.intercept",
  "router1": "",
  "router2": ""
}
```

There is an example over on [GONE-intercept](https://github.com/David-Antunes/gone-intercept) leveraging this type of operation.