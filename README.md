## DockerWatcher

**DockerWatcher** is a monitoring solution written in Go, designed to track Docker events and notify system administrators in real-time. By leveraging Docker's event system, DockerWatcher helps you stay informed about crucial changes such as container creation or termination, ensuring you are always aware of the state of your Docker environment.

### Key Features

- **Event Monitoring:** Capture and filter events from Docker's event system.
- **Multi-Platform Notifications:** Receive alerts on your preferred communication platforms, including Telegram (Supported), Discord (Soon), Skype (Soon), and more (Soon).
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
- [ ] Skype
- [ ] Slack
- [ ] Mattermost

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

### How to run

please see `.env.example` file. create another file like this and name it `.env`.

Also you can export these variables in the shell you use.

## Environment Variables

| Key                      | Description               | Optional/Required |
|--------------------------|---------------------------|-------------------|
| REDIS_URL                |                           | Required          |
| REDIS_QUEUE_WRITE_NAME   |                           | Required          |
| REDIS_QUEUE_READ_NAMES   |                           | Required          |
| ENABLE_WATCHER           | 0 (disabled), 1 (enabled) | Required          |
| ENABLE_NOTIFICATION      | 0 (disabled), 1 (enabled) | Required          |
| GRACEFUL_SHUTDOWN_TIMEOUT| Integer (seconds)         | Required          |
| EVENTS_FILTER            |                           | Optional          |
| NOTIFICATION_PLATFORM    | telegram                  | Required          |
| TELEGRAM_BOT_API_TOKEN   |                           | Required          |
| TELEGRAM_CHAT_ID         |                           | Required          |



`REDIS_QUEUE_READ_NAMES` and `EVENTS_FILTER` are comma seperated. watch `.env.example`.

for `EVENTS_FILTER`, see [this link](https://docs.docker.com/reference/cli/docker/system/events/#filter).

