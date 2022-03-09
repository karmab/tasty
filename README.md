# Tasty repository

Tasty is a CLI to handle with operators as well as you can integrate in your code to be used as a sdk to handle operators easily.

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
curl -s -L https://github.com/karmab/tasty/releases/latest/download/tasty-linux-amd64 > /usr/bin/tasty
chmod u+x /usr/bin/tasty
```

## Requirements

Kubeconfig environment variable must be set.
```
export KUBECONFIG=/path/to/kubeconfig
```

##  Running as kubelet and oc plugin

Run the following and you can then use `kubectl olm` or `oc olm`

```
tasty config --enable-as-plugin
```

## Use Tasty in your Code (Golang)

First of all, you need to import tasty:

```
import "github.com/karmab/tasty/pkg/operator"
```

Then, you need to create an operator using the constructor (empty to be filled after, or directly with the options):

```
o := operator.NewOperator()
```

or with options directly:

```
o := operator.NewOperatorWithOptions(name, source, defaultChannel, description, csv, namespace, crd, configExecFile, configExecPath string)
```

Now you could use the next functions to manage the operator:

```
o.SearchOperator()
o.GetList()
o.GetInfo()
o.SetConfiguration()
o.GetOperator()
o.InstallOperator()
```

## Installing operators from a pod

Check [job.yml.sample](job.yml.sample) as an example of a job that will install a given operator using a sa with cluster admin privileges

## Problems?

Open an issue!

Mc Fly!!!

karmab
