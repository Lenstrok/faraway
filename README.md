# Faraway Coding Challenge

## QuickStart

### Prerequisites

Before you begin, ensure you have the following installed on your system:

- [Go](https://golang.org/doc/install) (version 1.22 or higher)
- Docker (optional, for running app in containers)
- [golangci-lint](https://golangci-lint.run/welcome/install/) (for linters)

### Setup

1. **(Optional) Configure the environment**

   Edit the `env.example` file in the root dir of the project to set up your configuration options.

   | ENV                 | requirements                          | description                                                                |
      |---------------------|---------------------------------------|----------------------------------------------------------------------------|
   | `CLIENT_NUMBER_OF_QUOTES` | default:`5`                           | Number of requests from client to server                                   |
   | `QUOTES_FILE_PATH` | default:`internal/config/quotes.json` | Need to be updated for Docker run + you should mount file to the container |
   | `SERVER_PORT`       | default:`8080`                        | the port on which the server will be running                               |
   | `SERVER_CONN_TIMEOUT`    | default:`10s`                         |                                                                            |
   | `POW_TOKEN_SIZE`           | default:`16`                          |
   | `POW_NONCE_SIZE`       | default:`8`                           |
   | `POW_COMPLEXITY`           | default:`10`                          | Max target bits                                                            |

2. **Run Server Component**

   ```bash
    make run-server-app
    ```
    - Server works until `syscall.SIGTERM` will be received.
    - In case of Docker run we need to mount `quotes.json` file to the container and update env `QUOTES_FILE_PATH`.

3. **Start Client Component**

   ```bash
    make run-client-app
    ```
    - Client makes `CLIENT_NUMBER_OF_QUOTES` requests and stops.


4. **Optional**

   Run tests & linters
   ```bash
     make run-linters
     make run-unit-tests
   ```
   It would be good to add more unittest and benchmarks.
   Also, it's good to add End-to-end tests.

## Choice of the POW algorithm explanation

I chose a hash-based PoW algorithm, similar to the well-known Hashcash mechanism, for the following reasons:

- Simplicity and Effectiveness: A hash-based PoW algorithm is straightforward to implement, widely adopted and proven. 
- Adjustable Difficulty: The difficulty level can be easily adjusted by changing the number of leading zeroes required in the hash.
- Low Resource Consumption for the Server: The server only needs to verify the solution by computing a single hash

More complex PoW algorithms exist (e.g., memory-bound PoW or ASIC-resistant algorithms), but they tend to introduce significant overhead both in terms of complexity and resource usage. For this use case—protecting a TCP server from DDoS attacks—a hash-based PoW algorithm offers a good balance between security, simplicity, and performance.
