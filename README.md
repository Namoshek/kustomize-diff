# kustomize-diff

A go utility which allows to generate the diff between two Kustomization directories. This is for example useful for pull request reviews when the result of the Kustomization should be reviewed and not only the changes in the sources.

## Install

The latest version can be downloaded from the [releases](https://github.com/Namoshek/kustomize-diff/releases).
Make sure the downloaded binary is executable (`chmod -x`) and place it somewhere in the `$PATH` for easy access.

To download and install in one go (make sure to select the correct binary for your system), use:

```sh
wget https://github.com/Namoshek/kustomize-diff/releases/download/v0.3.0/kustomize-diff-v0.3.0-linux-amd64.tar.gz \
  && tar -xzvf kustomize-diff-v0.3.0-linux-amd64.tar.gz \
  && chmod a+x kustomize-diff \
  && sudo mv kustomize-diff /usr/local/bin/kustomize-diff
```

## Usage

Running `kustomize-diff` is as simple as:

```sh
$> kustomize-diff inline ./old-version/overlays/dev ./new-version/overlays/dev
```

which will give you a result (obviously depending on the differences) like:

````sh
```diff
 apiVersion: apps/v1
 kind: Deployment
 metadata:
   name: my-app
   namespace: my-namespace
 spec:
-  replicas: 1
+  replicas: 2
   selector:
     matchLabels:
       app: my-app
     template:
       metadata:
         labels:
           app: my-app
       spec:
         containers:
         - name: app
           image: my-app:latest
           ports:
-          - name: http
+          - name: https
             protocol: TCP
-            containerPort: 8080
+            containerPort: 8081
```
```diff
 apiVersion: v1
 kind: Service
 metadata:
   name: my-app
   namespace: my-namespace
 spec:
   selector:
     app: my-app
   type: ClusterIP
   ports:
-  - name: http
-    port: 8080
+  - name: https
+    port: 8081
     protocol: TCP
-    targetPort: http
+    targetPort: https
```
````

By default, `kustomize-diff` will use the `kustomize` binary from the `$PATH` to create the Kustomization of the given directories. By providing the `--kustomize-executable=<path>` option, a custom Kustomize executable may be used instead.

In case of success, the command will exit with the exit code `0`. Otherwise, an exit code `>0` will be returned.

### Diff for Pull Request Review

To use this utility in a pull request pipeline, it is recommended to checkout the source repository two times, once for the pull request target branch and once for the pull request source branch. The output of `kustomize-diff` can then be posted as pull request comment for review, for example.

#### Azure DevOps

To simplify posting the diff as comment(s) on a pull request, `kustomize-diff` provides a second command called `azuredevops` which accepts a few parameters that make this a breeze:

```sh
kustomize-diff azuredevops \
  --organization <organization> \
  --project <project> \
  --repository-id <repository-name-or-id> \
  --pull-request-id <pull-request-id> \
  --personal-access-token <pat> \
  --hide-diff-in-spoiler \
  --comment-per-resource \
  <pathToOldKustomization> <pathToNewKustomization>
```

## Development

### Run locally

The application can be run locally using:

```sh
go run main.go <arguments> [flags]
```

### Run Tests

To run the tests, use:

```sh
go test -v --cover ./...
```

### Build

To build the application, use:

```sh
env GOOS=<os> GOARCH=<arch> go build -o bin/kustomize-diff main.go
```

where `GOOS`/`GOARCH` are from `go tool dist list`, e.g. `GOOS=linux GOARCH=amd64` or `GOOS=windows GOARCH=amd64`.

## License

`Namoshek/kustomize-diff` is open-sourced software licensed under the [MIT license](LICENSE).
