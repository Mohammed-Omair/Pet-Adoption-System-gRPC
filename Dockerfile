# Step 1: Use the official Go image as the base image
FROM golang:1.23-alpine

# Step 2: Set the working directory inside the container
WORKDIR /app

# Step 3: Copy the Go modules files (if any) and the go.mod/go.sum first to the container
COPY go.mod go.sum ./

# Step 4: Download dependencies
#RUN go mod download

# Step 5: Copy the remaining files (like source code) to the container
COPY . .

# Step 6: Build the Go application
#RUN go build -o pet_adoption_server pet_adoption_server.go

# Step 7: Expose the port on which the server runs
EXPOSE 50051

# Step 8: Command to run the server when the container starts
#ENTRYPOINT [ "executable" ] ["runtime.sh"]
ENTRYPOINT ["/app/runtime.sh"]
#ENTRYPOINT ["sleep 10"]
