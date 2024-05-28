## DockerWatcher

**DockerWatcher** is a monitoring solution written in Go, designed to track Docker events and notify system administrators in real-time. By leveraging Docker's event system, DockerWatcher helps you stay informed about crucial changes such as container creation or termination, ensuring you are always aware of the state of your Docker environment.

### Key Features

- **Event Monitoring:** Capture and filter events from Docker's event system.
- **Multi-Platform Notifications:** Receive alerts on your preferred communication platforms, including Telegram (Supported), Mattermost (Supported), Generic webhook (Soon) and more (Soon).
- **Scalable Architecture:** Designed with a modular approach to efficiently handle and process event data.
- **Written in Go:** Utilizes the powerful and efficient Go programming language for high performance and reliability.

### Project Structure

DockerWatcher consists of two main services:

1. **Watcher Service:**
    - Collects and filters data from Docker's event system.
    - Pushes the relevant event data to Redis for further processing.
    
2. **Notification Service:**
    - Retrieves filtered event data from Redis.
    - Sends notifications to users via the configured communication platforms.

### Running Modes

DockerWatcher can be run in three different modes, providing flexibility based on your requirements:

1. **Watcher Mode:**
    - Collects data from Docker's event system and pushes it to Redis.
    - Ideal for distributed environments where multiple instances collect data from different Docker engines.

2. **Notification Mode:**
    - Retrieves data from Redis and sends notifications to users.
    - Useful for centralizing the notification logic in a single service.

3. **Watcher and Notification Mode:**
    - Combines both watcher and notification functionalities.
    - Collects data and sends notifications in one go, simplifying the deployment.

By using lists in Redis, DockerWatcher creates queues of messages, ensuring efficient and orderly processing of events.

### Notification Platforms

- [x] Telegram
- [x] Generic webhook
- [x] Mattermost


### Motivation

DockerWatcher was created to address the need for a reliable and efficient way to monitor Docker environments. By providing timely notifications about critical events, DockerWatcher helps system administrators maintain better control and oversight of their Docker containers, improving overall system reliability and responsiveness.

## Usage

### Build from Source Code

To build DockerWatcher from source, follow these steps:

1. **Clone the Repository:**
    ```sh
    git clone https://github.com/ravanbod/dockerwatcher.git
    cd dockerwatcher
    ```

2. **Build the Project:**
    ```sh
    go build -o dockerwatcher cmd/dockerwatcher/main.go
    ```

### Build with Docker

To build DockerWatcher with docker, follow these steps:

1. **Clone the Repository:**
    ```sh
    git clone https://github.com/ravanbod/dockerwatcher.git
    cd dockerwatcher
    ```

2. **Build the Project:**
    ```sh
    docker build -t dockerwatcher .
    ```

## How to run

### Environment Variables

please see `.env.example` file. create another file like this and name it `.env`.

Also you can export these variables in the shell you use.

| Key                      | Description               | Optional/Required |
|--------------------------|---------------------------|-------------------|
| REDIS_URL                |                           | Required                 |
| REDIS_QUEUE_WRITE_NAME   |                           | Required                 |
| REDIS_QUEUE_READ_NAMES   |                           | Required                 |
| ENABLE_WATCHER           | 0 (disabled), 1 (enabled) | Required                 |
| ENABLE_NOTIFICATION      | 0 (disabled), 1 (enabled) | Required                 |
| GRACEFUL_SHUTDOWN_TIMEOUT| Integer (seconds)         | Required                 |
| EVENTS_FILTER            |                           | Optional                 |
| NOTIFICATION_PLATFORM    |generic,telegram,mattermost| Required                 |
| TELEGRAM_BOT_API_TOKEN   |                           | Required if np=telegram  |
| TELEGRAM_CHAT_ID         |                           | Required if np=telegram  |
| GENERIC_NOTIFICATION_URL |                           | Required if np=generic   |
| MATTERMOST_HOST          |                           | Required if np=mattermost|
| MATTERMOST_BEARER_AUTH   |                           | Required if np=mattermost|
| MATTERMOST_CHANNEL_ID    |                           | Required if np=mattermost|

*np = NOTIFICATION_PLATFORM

for example, if you want to set NOTIFICATION_PLATFORM=telegram, TELEGRAM_BOT_API_TOKEN and TELEGRAM_CHAT_ID are necessary.

`REDIS_QUEUE_READ_NAMES` and `EVENTS_FILTER` are comma seperated. watch `.env.example`.

for `EVENTS_FILTER`, see [this link](https://docs.docker.com/reference/cli/docker/system/events/#filter).

### Run without docker

To run the project without docker, You can simply run this command (after Build from source code).

```
./dockerwatcher
```

### Run with docker

To run the project with docker, You can simply run this command (after Build with docker).

```
docker run --name dockerwatcher -it -v $(pwd)/.env:/app/.env -v /var/run/docker.sock:/var/run/docker.sock --network dockerwatcher dockerwatcher
```

A redis has to be in `dockerwatcher` network.

### Run with docker-compose (Recommended)

If you have an existing redis, you can use `docker-compose.yml`. else you can use `docker-compose-full.yml` that has redis and DockerWatcher.

#### Run multiple instances

It is recommended that launch a watcher service in every docker engine and use a shared Redis between them. And launch a Notification instance that sends Redis messages to your notification platform.

## Messages
You will get the messages in the markdown format. this is an example of dying mysql container.
```
# Docker Event 

 ## Event Details 

- **Type**: `container`
- **Action**: `die`
- **Scope**: `local`
- **Time**: `1716908784`
- **TimeNano**: `1716908784342468156`
## Actor 
- **Actor.ID**: `ab1bec9756b03eb3c12c42fe08496a94d73b037c8e60e41344ae869fd41e38cb`
  - **name**: `mysql`
  - **execDuration**: `5`
  - **exitCode**: `0`
  - **image**: `mysql`
```
