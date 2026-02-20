# FastPay Algeria 🇩🇿 💳

**FastPay** is a high-performance digital payment backend built with **Go (Golang)**, specifically designed for the Algerian market. It enables seamless peer-to-peer transfers, parent-child wallet management, and merchant payment processing.

## 🎯 Objective
To provide a secure, scalable, and ultra-fast financial ecosystem that supports:
* **Digital Wallets:** Unique 20-digit identifiers for every user.
* **Family Banking:** Parents can create and manage digital cards/wallets for their children.
* **Merchant Integration:** Support for E-commerce and physical stores via QR/NFC logic.
* **Security:** High-fidelity transaction logging and encrypted financial data.

## 🛠 Technology Stack
* **Language:** Go (1.21+) - chosen for its superior concurrency and speed.
* **Framework:** [Gin Gonic](https://github.com/gin-gonic/gin) (HTTP Web Framework).
* **Database:** PostgreSQL (Relational DB for ACID compliance).
* **ORM:** GORM (The fantastic ORM library for Golang).
* **Auth:** JWT (JSON Web Tokens) with Bcrypt password hashing.

## 📂 Project Structure
```text
├── cmd/api/            # Main entry point
├── internal/
│   ├── api/            # Handlers & Routes
│   ├── core/           # Business Logic & Domain Models
│   └── repository/     # Database Operations
├── pkg/                # Shared utilities (ID Generation, Security)
└── db/migrations/      # Database Schema files 