# Simple Golang sample using kafka and websocket

Simple signing service which involves communication between the web interface and signing service over Kafka.

Implementation flow
* User type some message in text field in the browser and click Submit button
* The message is submitted to REST API server which push SignEvent event to Kavka
* Signer Service listen for SignEvent event, sign received message, and push SignedEvent event with signed message. For signing you can use any algorithm you like, for example XOR.
* REST API listen for SignedEvent event, and stream signed message back to the browser over socket
* Browser listens for socket and shows received data. 


## Setup

before running application we need to start docker environment:
```shell
make start
```

then we need to build app:
```shell
make build
```

and finally run app:
```shell
./server
```

there is also possibility to build and deploy docker image:
```shell
make image
```

and to push it on docker hub:
```shell
make push
```