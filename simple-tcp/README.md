# Simple TCP Server

A simple game server for testing TCP based server on Carrier.

## Server

Starts a server on port `7654` by default. Can be overwritten by `PORT` env var or `port` flag.

When it receives a legal "\<Command\> TRUE" with a newline, it will send back "ACK: \<Command\>" as an echo. 

If it receives the text "EXIT", then it will `sys.Exit(0)`

## Run simple-tcp as a GameServer

```shell
# kubectl apply -f gameserver.yaml             
gameserver.carrier.ocgi.dev/simple-tcp-example created

# kubectl get pods
simple-tcp-example                       2/2     Running                 0          4s


# kubectl exec -it simple-tcp-example -c server -- /bin/bash
[root@simple-tcp-example /]# netstat -lntp
Active Internet connections (only servers)
Proto Recv-Q Send-Q Local Address           Foreign Address         State       PID/Program name    
tcp        0      0 127.0.0.1:9020          0.0.0.0:*               LISTEN      -                   
tcp        0      0 127.0.0.1:9021          0.0.0.0:*               LISTEN      -                   
tcp6       0      0 :::7654                 :::*                    LISTEN      1/simple-tcp        
tcp6       0      0 :::8080                 :::*                    LISTEN      -      
[root@simple-tcp-example /]# telnet 127.0.0.1 7654
Trying 127.0.0.1...
Connected to 127.0.0.1.
Escape character is '^]'.
FILLED TRUE
ACK: FILLED
INVALID COMMAND
Invalid command: INVALID
...
```

The `GameServer` condition will change:

```shell
# kubectl get gs/simple-tcp-example -o yaml
apiVersion: carrier.ocgi.dev/v1alpha1
kind: GameServer
...
status:
  conditions:
  - lastProbeTime: "2021-03-30T06:54:30Z"
    lastTransitionTime: "2021-03-30T06:54:30Z"
    status: "True"
    type: carrier.ocgi.dev/ready
  - lastProbeTime: "2021-03-30T06:55:58Z"
    lastTransitionTime: "2021-03-30T06:55:58Z"
    status: "True"
    type: carrier.ocgi.dev/filled
  state: Running
```

## Run simple-tcp as a Squad

```shell
# kubectl apply -f ./squad.yaml 
squad.carrier.ocgi.dev/squad-example created

# kubectl get sqd 
NAME              SCHEDULING       DESIRED   CURRENT   UP-TO-DATE   READY   AGE
squad-example     MostAllocated    2         2         2            0       5s

# kubectl get pods        
NAME                             READY   STATUS    RESTARTS   AGE
squad-example-66bcd47554-9vlzg     2/2     Running             0          53s
squad-example-66bcd47554-mrnd9     2/2     Running             0          53s
```

* Scale Up

```shell
# kubectl scale sqd squad-example --replicas=3
squad.carrier.ocgi.dev/squad-example scaled
# kubectl get sqd
NAME            SCHEDULING   DESIRED   CURRENT   UP-TO-DATE   READY   AGE
squad-example                3         3         3            3       6m
# kubectl get pods
NAME                               READY   STATUS              RESTARTS   AGE
squad-example-66bcd47554-9vlzg     2/2     Running             0          2m32s
squad-example-66bcd47554-mrnd9     2/2     Running             0          2m32s
squad-example-66bcd47554-tqlqh     2/2     Running             0          73s
```

* Scale Down

```shell
# kubectl scale sqd squad-example --replicas=2

## The number of replicas of squad will not decrease immediately because of the `deletableGates`
# kubectl get sqd 
NAME            SCHEDULING       DESIRED   CURRENT   UP-TO-DATE   READY   AGE
squad-example   MostAllocated    2         3         3            3       4m19s

## after 30 seconds
# kubectl get sqd
NAME            SCHEDULING       DESIRED   CURRENT   UP-TO-DATE   READY   AGE
squad-example   MostAllocated    2         2         2            2       6m17s
```
