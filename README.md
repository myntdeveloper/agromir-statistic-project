# Agromir Telemetry & Boiler Monitoring

**Production industrial IoT monitoring platform built for Agromir.** Ingests real-time JSON telemetry from Siemens S7 PLCs via HTTP, tracks boiler operational metrics across multiple warehouses, and provides detailed diagnostics and historical analytics.

[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![React](https://img.shields.io/badge/React-18-61DAFB?logo=react&logoColor=111827)](https://react.dev/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5-3178C6?logo=typescript&logoColor=white)](https://www.typescriptlang.org/)
[![Vite](https://img.shields.io/badge/Vite-6-646CFF?logo=vite&logoColor=white)](https://vitejs.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](#license)

## Key Features

- **PLC Telemetry Ingestion:** High-throughput Go API receiving real-time JSON payloads directly from Siemens S7 industrial controllers.
- **Multi-Warehouse Monitoring:** Live overview and drill-down diagnostics for boilers and cisterns across production facilities.
- **Historical Analytics:** Interactive time-series charts for tracking equipment performance, temperature, and pressure trends.
  
## Tech Stack
- **Backend:** Go (Gin), REST API, JSON Telemetry Ingestion
- **Hardware / IIoT:** Siemens S7 PLCs (HTTP / SCL)
- **Frontend:** React 18, TypeScript, Vite, Tailwind CSS
- **Database:** PostgreSQL 15
- **Infrastructure:** Docker Compose, Nginx

## Directory Structure

```text
agromir-statistic/
├── services/
│   ├── api/               # Backend (Golang, Gin)
│   └── frontend/          # Frontend (React + Vite)
├── infra/                 # Docker Compose, Nginx, infrastructure config
└── Readme.md              # Main English README
```

## Screenshots

### 1) Warehouse selection

<img width="50%" height="856" alt="image" src="https://github.com/user-attachments/assets/087bf2ab-92fa-4a04-9ad5-878635b205e4" />

### 2) Boiler selection

<img width="50%" height="856" alt="image" src="https://github.com/user-attachments/assets/da00fe68-92db-4a2e-96c9-81b1bfbb21a7" />


### 3) Boiler details

<img width="50%" height="856" alt="image" src="https://github.com/user-attachments/assets/26745a38-4dd4-430e-8019-e825b1106058" />


### 4) Telemetry & Analytics Charts 

<img width="50%" height="1151" alt="image" src="https://github.com/user-attachments/assets/1561798b-5805-4f1a-b0b2-924646d0b01d" />

---

Built with ❤️ by **mynt**
