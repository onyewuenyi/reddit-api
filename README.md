To set up a local Go development environment using a cloned repository inside a Docker container, you can follow these steps:

1. Install Docker: If you haven't already, install Docker on your machine by following the official Docker installation guide for your operating system.

2. Clone the repository: Clone the Go repository that you want to work on to your local machine using Git. For example, if the repository is hosted on GitHub, you can run the following command in your terminal:
   ```
   git clone <repository_url>
   ```

3. Create a Dockerfile: Inside the cloned repository directory, create a file named `Dockerfile` (without any file extension) with the following content:
   ```
    # Use a base image with Go preinstalled
    FROM golang:1.16-alpine

    # Set the working directory inside the container
    WORKDIR /app

    # Copy the source code and go.mod/go.sum files into the container
    COPY go.mod go.sum ./

    # Download the Go module dependencies
    RUN go mod download

    # Copy the rest of the source code into the container
    COPY . .

    # Build the Go application
    RUN go build -o my-go-app

    # # Set the command to run when the container starts
    # TODO find out if this cmd is needed or nah
    # CMD ["./my-go-app"]

   ```

   The Dockerfile uses the official `golang` base image, sets the working directory to `/app`, copies the contents of the current directory (cloned repository) to the container's `/app` directory, runs `go mod download` to download the project's dependencies, and finally specifies the command to run your Go application.

4. Build the Docker image: Open a terminal, navigate to the directory containing the `Dockerfile`, and run the following command to build the Docker image:
   ```
   docker build -t my-go-app .
   ```

   This command will build a Docker image with the tag `my-go-app`. Make sure to include the dot at the end, as it indicates the current directory as the build context.

5. Run the Docker container: Once the Docker image is built, you can run a container based on that image using the following command:
   ```
   winpty docker run -it --rm -p 8080:8080 -v "$(pwd)":/app my-go-app

   ```

   This command runs a Docker container interactively (`-it`), removes it after it stops (`--rm`), mounts the current directory (`$(pwd)`) to the container's `/app` directory, and uses the `my-go-app` image.

6. Develop inside the Docker container: Now you can start developing inside the Docker container. Any changes you make to the local directory will be reflected inside the container, allowing you to build and run your Go application.

That's it! You now have a local Go development environment set up inside a Docker container using the cloned repository. You can modify the Dockerfile as needed depending on your specific requirements, such as specifying additional dependencies or build configurations.