# ten-golang

Go implementation of TEN with MCTS AI and a simple Cloud Function endpoint.

It now reached performance parity with its [Java counterpart](https://github.com/SamChou19815/TEN)
on my local machine. There is not much performance difference between the Java and Go version, even
if the Go program will be compiled to native code. This is probably because Java has a highly
optimized garbage collector, since MCTS algorithm needs to create deep trees.

## Getting Started

### Running Locally

```bash
go run developersam.com/ten-golang/main
```

### Deploy to Google Cloud Functions

```bash
# We need more memory to create MCTS simulation trees.
gcloud functions deploy HandleTenAIMoveRequest --runtime go111 --trigger-http --memory=2048MB
```
