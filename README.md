# SAP Event Broker Take Home Problem

The following take home problem is to help us learn a a bit more about your coding style and creative problem solving. In
general we prefer this approach rather than forcing candidates to code during the interview. It is purposely very open ended
and you can tackle the problem in a wide range of ways.

## Guidelines

- Please do not spend more than 2 to 3 hours on the problem. Your time is valuable and there are no bonus points for candidates who do more than that.

- Think of this like an open source project. Create a repo on Github, use git for source control, and use README.md to document what you built for the newcomer to your project.

- Our team builds systems engineered to run in production. Given this, please organize, design, and document your solution as if you were going to put into production. We completely understand this might mean you can't do as much in the time budget. Be biased for production-ready over features.

- Document tradeoffs, the rationale behind your technical choices, or things you would do or do differently if you were able to spend more time on the project or do it again.

- Use whatever languages, libraries, tools etc. you feel like.

## The Problem

We the North American team who are split across Boulder and Montreal often enjoy spending time in the outdoors. Some of us would like to bike while some of us are happy just walking, some of us like to challenge ourselves with difficult hikes while some of us would like to ease in slowly.

The good part is that the Boulder county releases data about the trails as open data. Your assignment is to help us find a trail that suits our requirement.

This is a freeform assignment. You can write a web API that returns a set of trails. You can write a web frontend that visualizes the trails. We also spend a lot of time in the shell, so a CLI that gives us a couple of options would be great.

The only requirement for the assignment is that it allow us to filter the trails by at least 2 criteria. You can choose to pivot the choices based on any of the options in there - bike trail vs walking trail, with the option of fishing vs not etc.

Feel free to tackle this problem in a way that demonstrates your expertise of an area or takes you out of your comfort zone.

Boulder's trail head data is [located here](https://opendata-bouldercounty.hub.arcgis.com/datasets/3a950053bbef46c6a3c2abe3aceee3de_0/explore) you can download a CSV of the data from that site. We have included a copy of the CSV in this repo as well.

# The Solution: Trail Finder App

This project is a trail-finding application designed to manage and view trail data. It is built using the Go programming language, PostgreSQL for the database, and runs within a Docker container, making it deployable on Kubernetes clusters. The application is hosted on AWS EKS (Elastic Kubernetes Service).

## Features

- Manage trail information such as amenities, difficulty levels, fees, and more.
- Load default data from a CSV file.
- API endpoints for retrieving and managing trails.
- Deployable as a Docker container on Kubernetes.

## Live Demo

The application is accessible at:
[Trail Finder Live Demo](http://abf9a26ceb9a74f4ab5af09193ca15fb-1756722096.us-east-1.elb.amazonaws.com/trails)

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Kubernetes (kubectl)](https://kubernetes.io/docs/tasks/tools/)
- [eksctl](https://docs.aws.amazon.com/eks/latest/userguide/eksctl.html)
- [AWS CLI](https://aws.amazon.com/cli/)
- [Go 1.23.2+](https://go.dev/doc/install)
- PostgreSQL (local or remote)

## Local Setup

### 1. Clone the Repository

\`\`\`bash
git clone <repository_url>
cd sap-eb-take-home-problem
\`\`\`

### 2. Configure Environment Variables

Create a `.env` file in the root directory and add:

\`\`\`makefile
DB_CONN_STRING="postgresql://<username>:<password>@localhost:5432/trails_db"
\`\`\`

### 3. Install Dependencies

\`\`\`bash
go mod download
\`\`\`

### 4. Run Migrations

To run the database migrations:

\`\`\`bash
go run main.go migrate
\`\`\`

This will create the required database tables and columns.

### 5. Run the Application

To run the server locally:

\`\`\`bash
go run main.go --server
\`\`\`

The server will be available at [http://localhost:8080](http://localhost:8080).

### 6. Use the CLI

To use the CLI tool, run:

\`\`\`bash
go build -o trail-cli cli.go
./trail-cli <command>
\`\`\`

**Available CLI commands:**

- \`migrate\` - Run database migrations.
- \`loadcsv\` - Load trail data from a CSV file.

**Example:**

\`\`\`bash
./trail-cli loadcsv --file=BoulderTrailHeads.csv
\`\`\`

### 7. Testing

To run the test suite:

\`\`\`bash
go test ./...
\`\`\`

## Docker Setup

### 1. Build the Docker Image

\`\`\`bash
docker build -t trail-finder:latest .
\`\`\`

### 2. Run the Docker Container

\`\`\`bash
docker run -p 8080:8080 -e DB_CONN_STRING="postgresql://<username>:<password>@<db_host>:5432/trails_db" trail-finder:latest
\`\`\`

### 3. Push to AWS ECR

Authenticate Docker to AWS ECR:

\`\`\`bash
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <your-aws-account-id>.dkr.ecr.us-east-1.amazonaws.com
\`\`\`

Tag and Push the image:

\`\`\`bash
docker tag trail-finder:latest <your-aws-account-id>.dkr.ecr.us-east-1.amazonaws.com/trail-finder:latest
docker push <your-aws-account-id>.dkr.ecr.us-east-1.amazonaws.com/trail-finder:latest
\`\`\`

## Kubernetes Deployment

### 1. Deploy PostgreSQL

Apply the PostgreSQL deployment:

\`\`\`bash
kubectl apply -f k8s/postgres-deployment.yaml
kubectl apply -f k8s/postgres-service.yaml
\`\`\`

### 2. Deploy ConfigMap and Secret

Apply the ConfigMap and Secret for environment variables:

\`\`\`bash
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secret.yaml
\`\`\`

### 3. Deploy the Trail Finder App

\`\`\`bash
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
\`\`\`

### 4. Check Application Status

\`\`\`bash
kubectl get pods
kubectl get svc
\`\`\`

To check logs of a running pod:

\`\`\`bash
kubectl logs <pod_name>
\`\`\`

## Interacting with the API

### 1. Get Trails

Retrieve a list of trails:

\`\`\`bash
curl -X GET "http://localhost:8080/trails?page=1&limit=5"
\`\`\`

### 2. Load Trails from CSV (via API)

\`\`\`bash
curl -X POST -H "Content-Type: application/json" -d '{"file_path": "./BoulderTrailHeads.csv"}' "http://localhost:8080/loadcsv"
\`\`\`

### 3. Filter Trails

\`\`\`bash
curl -X GET "http://localhost:8080/trails?restrooms=yes&fishing=yes"
\`\`\`

## Project Structure

\`\`\`bash
sap-eb-take-home-problem/
│  
├── cmd/ # CLI tool
├── db/ # Database-related code
├── handlers/ # API handlers
├── k8s/ # Kubernetes deployment files
├── migrations/ # Database migration files
├── models/ # Models for the application
├── tests/ # Test files
├── BoulderTrailHeads.csv # Default data
├── cli.go # CLI entrypoint
├── docker-compose.yml # Docker Compose file
├── Dockerfile # Dockerfile for building the app
├── go.mod # Go module file
├── main.go # Main application file
└── README.md # This documentation
\`\`\`

## Technologies Used

- **Go**: Programming language for building the server and CLI.
- **PostgreSQL**: Relational database for storing trail data.
- **Docker**: Containerization for app deployment.
- **Kubernetes**: Orchestration for scaling and managing the app.
- **AWS ECR & EKS**: Hosting and deployment of containers on the cloud.

## Troubleshooting

- **Database Connection Issues**: Ensure the `DB_CONN_STRING` is correct and the database is accessible.
- **Kubernetes Pods Not Starting**: Check logs using `kubectl logs <pod_name>` to diagnose errors.
- **CSV File Not Found**: Ensure the file is properly mounted in the container using ConfigMap.

## Author

Pouria Roostaei
