# Deprecation Notice
> [!IMPORTANT]  
> This repository is no longer maintained. 
> 
> The content has moved to [github.com/loft-sh/vcluster-sdk examples](https://github.com/loft-sh/vcluster-sdk/tree/main/examples)
> 
## vCluster Plugin Example

This example plugin syncs car crds from the vcluster into the host cluster.

For more information how to develop plugins in vcluster and a complete walk through, please refer to the [official vcluster docs](https://www.vcluster.com/docs/plugins/overview).

### Using the Plugin

To use the plugin, create a new vcluster with the `plugin.yaml`:

```
# Apply cars crd in host cluster
kubectl apply -f https://raw.githubusercontent.com/loft-sh/vcluster-plugin-example/main/manifests/crds.yaml

# Use public plugin.yaml
vcluster create my-vcluster -n my-vcluster -f https://raw.githubusercontent.com/loft-sh/vcluster-plugin-example/main/plugin.yaml
```

This will create a new vcluster with the plugin installed. After that, wait for vcluster to start up and check:

```
# Create a car in the virtual cluster
vcluster connect my-vcluster -n my-vcluster -- kubectl apply -f manifests/audi.yaml

# Check if the car was synced to the host cluster
kubectl get car -n my-vcluster
```

### Building the Plugin
To just build the plugin image and push it to the registry, run:
```
# Build
docker build . -t my-repo/my-plugin:0.0.1

# Push
docker push my-repo/my-plugin:0.0.1
```

Then exchange the image in the `plugin.yaml`.

## Development

General vcluster plugin project structure:
```
.
├── go.mod              # Go module definition
├── go.sum
├── devspace.yaml       # Development environment definition
├── devspace_start.sh   # Development entrypoint script
├── Dockerfile          # Production Dockerfile 
├── main.go             # Go Entrypoint
├── plugin.yaml         # Plugin Helm Values
├── syncers/            # Plugin Syncers
└── manifests/          # Additional plugin resources
```

Before starting to develop, make sure you have installed the following tools on your computer:
- [docker](https://docs.docker.com/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/) with a valid kube context configured
- [helm](https://helm.sh/docs/intro/install/), which is used to deploy vcluster and the plugin
- [vcluster CLI](https://www.vcluster.com/docs/getting-started/setup) v0.6.0 or higher
- [DevSpace](https://devspace.sh/cli/docs/quickstart), which is used to spin up a development environment
- [Go](https://go.dev/dl/) programming language build tools

After successfully setting up the tools, start the development environment with:
```
devspace dev -n vcluster
```

After a while a terminal should show up with additional instructions. Enter the following command to start the plugin:
```
go build -mod vendor -o plugin main.go && /vcluster/syncer start
```

You can now change a file locally in your IDE and then restart the command in the terminal to apply the changes to the plugin.

Delete the development environment with:
```
devspace purge -n vcluster
```
