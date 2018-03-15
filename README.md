[![Go Report Card](https://goreportcard.com/badge/github.com/pharmer/flexvolumes)](https://goreportcard.com/report/github.com/pharmer/flexvolumes)
[![Build Status](https://travis-ci.org/pharmer/flexvolumes.svg?branch=master)](https://travis-ci.org/pharmer/flexvolumes)
[![codecov](https://codecov.io/gh/pharmer/flexvolumes/branch/master/graph/badge.svg)](https://codecov.io/gh/pharmer/flexvolumes)
[![Docker Pulls](https://img.shields.io/docker/pulls/pharmer/flexvolumes.svg)](https://hub.docker.com/r/pharmer/flexvolumes/)
[![Slack](https://slack.appscode.com/badge.svg)](https://slack.appscode.com)
[![Twitter](https://img.shields.io/twitter/follow/appscodehq.svg?style=social&logo=twitter&label=Follow)](https://twitter.com/intent/follow?screen_name=AppsCodeHQ)

# flexvolumes

This is a collection of Kubernetes FlexVolume plugins. Flexvolume is a GA feature from Kubernetes 1.8 release onwards.

### Build

```
go get github.com/pharmer/flexvolumes
```

### Install

See [here](hack/deploy). Fill in your DigitalOcean API key in `secret.yaml` and upload:

```
kubectl create -f secret.yaml
```

Add`--enable-controller-attach-detach=false` to `KUBELET_ARGS`.

Restart kubelet with `systemctl restart kubelet`.

After that, run:

```
kubectl create -f daemonset.yaml
```

### Usage

Create a volume on DigitalOcean if you haven't already done so, and find
its id. From the website it looks like the best way to do it is to inspect
element on the volumes page and look for a div with data-id="...", or use their
API, or if you're using terraform inspect the state. Fill in the id in pod.yaml.
Then:

```
kubectl create -f pod.yaml
```

---

**Pharmer binaries collects anonymous usage statistics to help us learn how the software is being used and how we can improve it. To disable stats collection, run the operator with the flag** `--analytics=false`.

---

## Support
We use Slack for public discussions. To chit chat with us or the rest of the community, join us in the [Kubernetes Slack team](https://kubernetes.slack.com/messages/C81LSKMPE/details/) channel `#pharmer`. To sign up, use our [Slack inviter](http://slack.kubernetes.io/).

To receive product announcements, please join our [mailing list](https://groups.google.com/forum/#!forum/pharmer) or follow us on [Twitter](https://twitter.com/AppsCodeHQ). Our mailing list is also used to share design docs shared via Google docs.

If you have found a bug with Pharmer or want to request for new features, please [file an issue](https://github.com/pharmer/pharmer/issues/new).
