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
curl -s -L https://github.com/karmab/tasty/releases/download/v0.4.0/tasty-linux-amd64 > /usr/bin/tasty
chmod u+x /usr/bin/tasty
```

##  Running as oc plugin

Run the following and you can then use `oc olm`

```
TASTYDIR=$(dirname $(which tasty))
ln -s $TASTYDIR/tasty $TASTYDIR/oc-olm
```

## Installing operators from a pod

Check [job.yml.sample](job.yml.sample) as an example of a job that will install a given operator using a sa with cluster admin privileges

## Problems?

Open an issue!

Mc Fly!!!

karmab
