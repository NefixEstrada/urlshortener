#
# URL Shortener, Néfix Estrada, 2018
# https://gitea.nefixestrada.com/nefix/urlshortener
#

#
# Build stage
#

# Use golang 1.11 as stage stage
FROM golang:1.11 as build

# Download the URL Shortener
RUN go get -d -v gitea.nefixestrada.com/nefix/urlshortener

# Move to the correct directory
WORKDIR /go/src/gitea.nefixestrada.com/nefix/urlshortener

# Download all the dependencies
RUN go get -d -v

# Compile the binary
RUN make

# Create the user
RUN adduser -D -g '' app

#
# Base stage
#

# Use alpine 3.8 as base
FROM alpine:3.8

# Copy the /etc/passwd (which contains the user 'app') from the build stage
COPY --from=build /etc/passwd /etc/passwd

# Copy the compiled binary from the build stage
COPY --from=build /go/src/gitea.nefixestrada.com/nefix/urlshortener/urlshortener /srv

# Use the 'app' user
USER app

# Expose the required port
EXPOSE 3000

# Run the service
CMD [ "/srv/urlshortener" ]
