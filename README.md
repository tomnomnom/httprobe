# httprobe

Take a list of domains and probe for working http and https servers.

## Install

```
▶ go get -u github.com/tomnomnom/httprobe
```

## Basic Usage

httprobe accepts line-delimited domains on `stdin`:

```
▶ cat recon/example/domains.txt
example.com
example.edu
example.net
▶ cat recon/example/domains.txt | httprobe
http://example.com
http://example.net
http://example.edu
https://example.com
https://example.edu
https://example.net
```

## Extra Probes

By default httprobe checks for HTTP on port 80 and HTTPS on port 443. You can add additional
probes with the `-p` flag by specifying a protocol and a comma-delimited port list pair:

```
▶ cat domains.txt | httprobe -p http:81,8080,8081 -p https:8443
```

If you want to probe some ports on both HTTP and HTTPS, then you can omit the protocol:

```
▶ cat domains.txt | httprobe -p 1234,4321,8088
```

You can probe a pre-defined list of ports that are commonly used for HTTP(S) by using the keywords `large` and `xlarge`:

```
▶ cat domains.txt | httprobe -p large
▶ cat domains.txt | httprobe -p https:large,1234 -p http:xlarge
```

## Concurrency

You can set the concurrency level with the `-c` flag:

```
▶ cat domains.txt | httprobe -c 50
```

## Timeout

You can change the timeout by using the `-t` flag and specifying a timeout in milliseconds:

```
▶ cat domains.txt | httprobe -t 20000
```

## Skipping Default Probes

If you don't want to probe for HTTP on port 80 or HTTPS on port 443, you can use the
`-s` flag. You'll need to specify the probes you do want using the `-p` flag:

```
▶ cat domains.txt | httprobe -s -p https:8443
```

## Prefer HTTPS

Sometimes you don't care about checking HTTP if HTTPS is working. You can do that with the `--prefer-https` flag:

```
▶ cat domains.txt | httprobe --prefer-https
```

## Docker

Build the docker container:

```
▶ docker build -t httprobe .
```

Run the container, passing the contents of a file into stdin of the process inside the container. `-i` is required to correctly map `stdin` into the container and to the `httprobe` binary.

```
▶ cat domains.txt | docker run -i httprobe <args>
```

