# kustomize-diff

A go utility which allows to generate the diff between two Kustomization directories. This is for example useful for pull request reviews when the result of the Kustomization should be reviewed and not only the changes in the sources.

## Usage

The latest version can be downloaded from the [releases](https://github.com/Namoshek/kustomize-diff/releases).

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
   namespace: foo
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
   namespace: foo
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

## License

`Namoshek/kustomize-diff` is open-sourced software licensed under the [MIT license](LICENSE).
