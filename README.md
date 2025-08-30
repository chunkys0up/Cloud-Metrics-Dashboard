# Cloud Metrics Dashboard
Lightweight cloud metrics dashboard that uses Typescript, Go, and Redis. Supports real time API request monitoring as well as aggregates data across multiple projects.

<img width="1506" height="818" alt="image" src="https://github.com/user-attachments/assets/524a23b6-2e9c-4348-b6ff-e74954dd20d2" />

## The dashboards collects and displays:
- Cpu, Memory, and Disk usage
- Rx and Tx Bytes rate
- Latency (ms)
- Total and failed requests

## Installation

### Prerequisites
Redis, Go and Node.js

### Clone Repo
```
git clone https://github.com/chunkys0up/Cloud-Metrics-Dashboard.git
cd Cloud-Metrics-Dashboard
```

### Install necessary packages for backend
```bash
cd server
go mod tidy
```

# Running the Project
Terminal 1
```bash
cd Frontend-dashboard
npm install
npm run dev
```

Terminal 2
Run the first line once, then anytime you want to run it again, just the 2nd one
```bash
go chmod +x ./script.sh 
./script.sh
```

Terminal 3 (or whatever project)
```
cd Test
go run TestLatency.go
```








