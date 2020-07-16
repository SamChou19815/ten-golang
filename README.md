# ten-golang

[![Build Status](https://github.com/SamChou19815/ten-golang/workflows/CI/badge.svg)](https://github.com/SamChou19815/ten-golang)

Go implementation of TEN with MCTS AI and a simple Cloud Function endpoint.

It now reached performance parity with its [Java counterpart](https://github.com/SamChou19815/ten-java)
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
gcloud functions deploy HandleTenAIMoveRequest --runtime go113 --trigger-http --memory=2048MB
```

## Additional Rules

The main rules are described [here](https://mathwithbaddrawings.com/2013/06/16/ultimate-tic-tac-toe).

To balance the game, I specified an additional rule that when there is same number of big squares
for black and white, white wins. It can compensate for the first-move advantage for black. With this
rule, the winning probability for black and white is 53:47. Without the rule, the ratio is 7:3.
