### captured_thoughts.md
SO for example, these are ways I plan on designing such

1. Collecting logs over http. For example, the application does some kind of async logging
to a central point, and at certain durations, a background thread sends the logs over http[s]

2. Another is via websockets, although we could make the http-longlived, websockets might prove less expensive 
as no reconnections are made


So we might have a whole http-server-thread running in the background[lol] or have a sep process?
Where the db logs events to a file, and the server-process talks w the ingestor. 
This way, the DB does DB work, and then the server-process just tails a file


We can call the server `cho`
We can call the ingestor `calatrava`


I'd want cho to be written in go, because the fileAPIs are just way more better than the java one
And then calatrava can be written in Spring


BUt why would you need to observe logs?
1. To know when your app is about to sink. Seen somewhere before that oncall engrs usually recv alarms that errorRate is at 5% or so
2. Measuring things like latency or performance? but I'm not sure
3. Suspicious activity? or knowing if the app is even live?


Hurdles
---
1. impl structured logging in the db. it usings log4j so it's okay at least
2. setting up springBoot and buildTools [since the underlying os is on nix]
3. schema evolution for the logs. [it seems necessary, but idk yet]



Experiment
---
Have two servers running **jkvs**. And just have requests blasted at them. They don't have to be in sync
neither do they need consistency among them. The plan here is to see if the log-ingestor/observer can 
pool from diff services


Another thing is that log-observers from my perspective, at this point is just for easy troubleshooting. I think 
it's supposed to get more powerful when merging logs from different distributed systems like a cluster?
Seems like logging to a file and sending it over http[s] but just glorified. And better with browser UI and UX Iguess
Almost feels like most of the work apart from extracting teh  data and parsing it on the other end would go 
to sorting and filtering and classifying?


For the protocol, I'd want something flexible and not very much rigid. something like gRPC? [haven't used before. Have
used RPC's and Protobufs seperately but never both at once, what does Google do in it?] mainly for the novelty, but that defeats 
the point of springboot

