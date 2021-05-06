# recover-geth-funds
Very simple cli to send all the funds from a geth KeyStore to an address.

```
COMMANDS:
   info     get info
   sendall  send all eth to address
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value
   --address value
   --help, -h       show help
   --version, -v    print the version
```

# Example

```
> go run main.go --address=0xbbb...bbb sendall

Sending to: 0xbbb...bbb
Keystore and account unlocked successfully: 0xaaa...aaa
Connection to web3 server opened url https://web3.url
Current balance: 0.361380906720109598 ETH
Nonce: 183
balance 0.10380807320336598 ETH
 tosend 0.10380802320336598 ETH (substracting the fees)
(y/N): y
operation confirmed
tx sent: 0xccccc...ccc
```
