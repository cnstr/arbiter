# Canister Arbiter
> Dead simple orchestration for Canister workers.

Canister Arbiter is a simple tool for orchestrating Canister workers.
These workers (also known as indexers) are responsible for building data that is consumed by the [Canister API](https://github.com/cnstr/api).

## Canister Peering
The current iteration of the Canister indexer runs as a single process that goes through a list of manifests and builds the data.
The issue is that this process can take a very long time to complete, making over 20,000 requests to various sources.
Measures have been implemented to avoid unnecessary work, but a full index can take almost an hour to complete.
On top of that, the indexer tends to trigger rate limits on some sources, leading to temporary IP bans.

To address these issues, we are splitting workloads across multiple volunteered peers.
Individuals have agreed to run a worker on their own infrastructure that gets directed by the orchestrator.
The plan is to eventually make the API, database, and orchestrator decentralized.

## Contributing
We'd appreciate anyone who is willing to host a peer for the Canister indexer.
This not only increases the reliability of the service, but also helps serve other regions faster.
In order to join the peering network, there are a few requirements that need to be met:

- Around `256MB` or less of RAM available at bursts every hour.
- Barely any CPU, it'll only occasionally burst up to half a thread.
- Relatively stable internet connection capable of opening PostgreSQL connections.
- You don't need to run the worker 24/7, but as much as possible is appreciated.
- Allow WebSocket connections to `i1.peer.canister.me`, `i2.peer.canister.me`, and `i3.peer.canister.me` on port `443`.

Canister's Peering Network is secured by a custom certificate authority to ensure that only approved peers can participate.
In order to join, you'll need to follow the steps below to generate a certificate signing request which you can send to us.
To generate a certificate signing request, you'll need `openssl` and the following `config.conf` file:

```conf
[req]
default_bits = 4096
prompt = no
encrypt_key = no
default_md = sha512
distinguished_name = req_distinguished_name

[req_distinguished_name]
C = US
ST = New York
O = Canister
OU = Indexer
CN = <name>
```

Remember to replace `<name>` with the name of the peer. And then run the following command:
```sh
openssl req -new -keyout <name>.key -out <name>.csr -config config.conf
```

Once you've generated the files, keep the `.key` file safe and send the `.csr` file to us.
Email [support@canister.me](mailto:support@canister.me) with the subject: `[PEER REQUEST] Indexer: <name>`.
Attach the `.csr` file, an explanation of your infrastructure and your relevance within the community.
We'll get back to you with a `.crt` file that you can use with the instructions below.

## Running a Peer
TODO
