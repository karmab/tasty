# Tasty repository

## Demo!

![](tasty.gif)

# Description

This repo provides a basic tool to interact with olm in a package manager way:

- list
- info
- install
- remove
- search

## Installing

```
curl -s -L https://github.com/karmab/tasty/releases/download/v0.3.0/tasty-linux-amd64 > /usr/bin/tasty
chmod u+x /usr/bin/tasty
```

##  Running as oc plugin

Run the following and you can then use `oc olm`

```
TASTYDIR=$(dirname $(which tasty))
ln -s $TASTYDIR/tasty $TASTYDIR/oc-olm
```

## Problems?

Open an issue!

Mc Fly!!!

karmab
