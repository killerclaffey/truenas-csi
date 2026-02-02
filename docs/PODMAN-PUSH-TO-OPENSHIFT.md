# Pushing Local Images to OpenShift Registry with Podman (Windows/PowerShell)

## Prerequisites
- Podman installed and machine running (`podman machine start`)
- `oc` CLI logged into your OpenShift cluster
- Image already built locally

## Procedure

### 1. Start Podman Machine (if not running)
```powershell
podman machine start
```

### 2. Get OpenShift Credentials
```powershell
$user = oc whoami
$token = oc whoami -t
```

### 3. Login to OpenShift Registry
```powershell
podman login default-route-openshift-image-registry.apps.okd.claffey.cloud -u $user -p $token --tls-verify=false
```

### 4. Tag Your Image for the Registry
```powershell
podman tag localhost/your-image:tag default-route-openshift-image-registry.apps.okd.claffey.cloud/namespace/your-image:tag
```

### 5. Push the Image
```powershell
podman push default-route-openshift-image-registry.apps.okd.claffey.cloud/namespace/your-image:tag --tls-verify=false
```

**Note:** Add `--format docker` if you encounter manifest format issues.

### 6. Verify in OpenShift
```powershell
oc get imagestream your-image -n namespace
```

### 7. Update Deployment to Use New Image
```powershell
oc set image deployment/your-deployment container-name=image-registry.openshift-image-registry.svc:5000/namespace/your-image:tag -n namespace
```

## Troubleshooting

### Podman Can't Connect
```powershell
podman machine stop
podman machine start
```

### Push Hangs or Fails
- Verify registry route: `oc get route -n openshift-image-registry`
- Check registry pods: `oc get pods -n openshift-image-registry`
- Try with `--format docker` flag

### Image Not Showing in ImageStream
After push, import explicitly:
```powershell
oc import-image your-image:tag --from=default-route-openshift-image-registry.apps.okd.claffey.cloud/namespace/your-image:tag -n namespace --confirm --insecure
```

## Example: TrueNAS CSI Driver

```powershell
# Build
cd C:\Users\rclaf\sync\truenas-csi
podman build -t truenas-csi-driver:v1.0.1 .

# Tag
podman tag truenas-csi-driver:v1.0.1 default-route-openshift-image-registry.apps.okd.claffey.cloud/truenas-csi/truenas-csi-driver:v1.0.1

# Login
$user = oc whoami
$token = oc whoami -t
podman login default-route-openshift-image-registry.apps.okd.claffey.cloud -u $user -p $token --tls-verify=false

# Push
podman push default-route-openshift-image-registry.apps.okd.claffey.cloud/truenas-csi/truenas-csi-driver:v1.0.1 --tls-verify=false

# Update deployment
oc set image deployment/truenas-csi-controller csi-controller=image-registry.openshift-image-registry.svc:5000/truenas-csi/truenas-csi-driver:v1.0.1 -n truenas-csi
oc set image daemonset/truenas-csi-node csi-node=image-registry.openshift-image-registry.svc:5000/truenas-csi/truenas-csi-driver:v1.0.1 -n truenas-csi
```
