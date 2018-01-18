[![Go Report Card](https://goreportcard.com/badge/github.com/pharmer/flexvolumes)](https://goreportcard.com/report/github.com/pharmer/flexvolumes)

# flexvolumes

This is a collection of Kubernetes FlexVolume plugins. So far I just have
DigitalOcean, and Packet is in progress. Since FlexVolumes are unstable, so too
is this. Use at your own risk. Contributions are welcome.

### Build

```
go get github.com/pharmer/flexvolumes
./hack/maky.py
```

### Install

Copy the plugin binary (e.g., `digitalocean`) to the Kubernetes volume plugin
directory:

```
mkdir -p /usr/libexec/kubernetes/kubelet-plugins/volume/exec/digitalocean/digitalocean
cp digitalocean /usr/libexec/kubernetes/kubelet-plugins/volume/exec/digitalocean/digitalocean
```

Note that CoreOS mounts `/usr` as read-only so instead you'll want to add
`--volume-plugin-dir=/etc/kubernetes/volumeplugins` to `KUBELET_ARGS` in
`/etc/kubernetes/kubelet.env` and put the plugins there instead.

Restart kubelet with `systemctl restart kubelet.service`.

### Usage

See `example`. Fill in your DigitalOcean API key in `secret.yaml` and upload:

```
kubectl create -f secret.yaml
```

Next, create a volume on DigitalOcean if you haven't already done so, and find
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
