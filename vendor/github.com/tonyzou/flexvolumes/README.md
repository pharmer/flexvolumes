# flexvolumes

This is a collection of Kubernetes FlexVolume plugins. So far I just have
DigitalOcean, and Packet is in progress. Since FlexVolumes are unstable, so too
is this. Use at your own risk. Contributions are welcome.

### Build

```
go get github.com/kardianos/govendor
govendor sync
go build ./provider/digitalocean
go build ./provider/packet
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
