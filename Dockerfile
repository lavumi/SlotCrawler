FROM alpine:latest

# Creates an app directory to hold your app’s source code
WORKDIR /

# Copies everything from your root directory into /app
ADD ./web ./web
COPY build/slot-crawler ./

# Tells Docker which network port your container listens on
EXPOSE 8081

# Specifies the executable command that runs when the container starts
CMD [ "/slot-crawler" ]