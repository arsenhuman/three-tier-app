# 🐳 Three-Tier Application Challenge with Docker (Practical Exam)

## Problem

You are tasked with deploying a **three-tier web application** using **Docker Compose**.

The stack consists of:

1. **Frontend** - a static web page served via Nginx  
2. **Backend** - a Go-based REST API  
3. **Database** - a MySQL instance with preloaded data  

#### Architecture
![](architecture.png)

#### Project Structure
```
/three-tier-app-challenge/
├── docker-compose.yaml        # Compose file defining all services
├── frontend/                  # Frontend static app
│   ├── Dockerfile             # Dockerfile for frontend
│   └── index.html             # Frontend HTML
└── backend/                   # Backend Go API
    ├── Dockerfile             # Dockerfile for backend
    ├── main.go                # Go source code
    ├── go.mod                 # Go modules
    └── go.sum                 # Go checksum file
```

**Your goal is to connect all components properly and display a message from the database on the frontend along with a visit counter.**


## Requirements

### **Database**

1. Use **MySQL** from the provided image `ghcr.io/hayk96/three-tier-app-challenge-mysql:latest`.
3. Database is already initialised with the following details and credentials. (to be used by **backend**): 

| Key          | Value         |
|--------------|---------------|
| username     | `student`     |
| password     | `student1234` |
| database name| `student_db`  |


### **Backend**

1. Build a **Go-based REST API** from `./backend` to serve:
   - JSON response: with the message and page visit counter from the database.
2. The backend uses environment variables to connect to the database. The variables are misconfigured in the codebase and must be set correctly in the `docker-compose.yaml`. Check the correct values in the **Database** section above.  
3. Expose the backend on **port 9000** to ensure it is reachable by the frontend (user browser).


3. The backend image **must be optimized** and **must not exceed 10 MB** in size.

### **Frontend**

1. Serve a static HTML page from `./frontend` using **Nginx**.  
2. The frontend must **only communicate with the backend**; direct access to the database is not allowed.  
3. Expose the frontend on **port 8080** and ensure it displays:
   - A **message** retrieved from the database via backend including **visit counter** that increments on every page refresh.

### **Networking Rules**

- Create two **user-defined networks**:
  - `frontend-network` → connects `frontend` ↔ `backend`
  - `backend-network` → connects `backend` ↔ `database`
- Assign the **backend container a static IP**: `172.20.0.10`. (see reference [here](https://docs.docker.com/reference/compose-file/services/#ipv4_address-ipv6_address), how to assign a static IP address to the container. 
- Subnets:
  - `frontend-network`: `X.X.X.X/X`
  - `backend-network`: `Y.Y.Y.Y/Y`

## Before You Submit

If you see that your frontend application correctly displays a message about completing the challenge and shows a visit 
counter that increments with each refresh, that's great. However, you need to consider that in case you destroy the stack, 
you will lose your data and current state. To prevent this, you need to enable persistence for your database. The database 
holds data in the `/var/lib/mysql` directory, so you need to mount a volume to that directory. Name the volume as `database-volume`.

## Task verification

To verify if database persistence is working as expected, please bring down the Docker Compose setup using `docker-compose down` 
and then recreate it. If, upon returning to the frontend page (`http://localhost:8080`), the visit counter continues incrementing from the previous 
value, you have successfully completed the challenge! 👏

## License
Distributed under the MIT License. See `LICENSE` for more information.

<!-- CONTACT -->
## Author and maintainer
**Hayk Davtyan | [@hayk96](https://github.com/hayk96)**# Three-Tier-App-Challenge-with-Docker
