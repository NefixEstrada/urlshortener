# URL Shortener

This is a simple URL shortener. You can visit it live at [https://short.nefixestrada.com](https://short.nefixestrada.com)

## Tech stack

The program is written in [Go](https://golang.org). For database, it's using [bbolt](https://github.com/etcd-io/bbolt), a fork of [bolt](https://github.com/boltdb/bolt). Also using [govalidator](https://github.com/asaskevich/govalidator) for validating the URL's and [go.rice](https://github.com/GeertJohan/go.rice) for embedding the static html files into the binary.

## How to run 

In order to run URL Shortener, you have two options:

### Inside a Docker [recommended]

You can easily run URL Shortener inside a Docker container. It's the recommended choice. You just need to download it from Docker Hub:

```sh
sudo docker pull nefix/urlshortener:1
sudo docker run -p 3000:3000 nefix/urlshortener
```

### As a standalone binary

You also can run it as a standalone binary. You need to execute the following commands:

```sh
git clone https://gitea.nefixestrada.com/nefix/urlshortener
cd urlshortener
make
```

This is going to generate a binary named `urlshortener`

## Examples

### Docker Compose

```yml
version: '3.2'
services:
  urlshortener_server:
    volumes:
      - type: volume
        source: urlshortener_data
        target: /data
        read_only: false
    ports:
      - target: 3000
        published: 8080
        protocol: tcp
        mode: host
    restart: always
    image: nefix/urlshortener:1

volumes:
  urlshortener_data:
```

## FAQ

- Why it doesn't have support for HTTPS?  
    + You are supposed to run it behind a proxy. If you have no idea what it is, you can check [here](https://en.wikipedia.org/wiki/Reverse_proxy) for more inforrmation.
- For any other question, you can contact me at [nefixestrada@gmail.com](mailto:nefixestrada@gmail.com)
