# kube-dash-backend

## How to run

```shell
$ openssl rand -hex 32 > secret.key
$ go build
$ ./kube-dash-backend -kubeconfig ~/.kube/config  # in case kubernetes controller is installed on your machine
```

### Docker

```shell
$ docker buildx build -t kube-dash-backend:latest .
$ docker run --rm -it -p 5000:5000 -v ./secret.key:/config/secret.key:ro -v ./kube-config.yml:/config/kube-config.yml:ro --name backend kube-dash-backend:latest
```

Assuming `secret.key` and `kube-config.yml` are in current working directory. Please adjust the mount points for your workflow.

## TODO

- [ ] change hard coded user credentials and their location
- [ ] make app suitable to run in [k8s itself](https://www.youtube.com/watch?v=NeV-jR_LssA)
- [ ] add more sophisticated JWT session control (now tokens are active for 2 hours)
- [x] more endpoints
- [x] better control over data sent to the user from those endpoints (process data)
- [x] make http error codes more uniform

## Swagger
To run swagger start the `kube-dash-backend` with `-swag -dev` flags. Then head to `http://localhost:5000/swagger` to see the API documentation. Disable both of those flags when running in production.

## API calls

> Slowly being moved to swagger, may not be accureate in the future

- `/api/v1/login` (POST) - login using hard coded user = 'john', pass = 'doe'. Will return JWT token which would need to be included in `Authorization: Bearer ...` header

- `/api/v1/accessible` (GET) - accessible endpoint, anyone can access  (only for testing)
- `/api/v1/restricted` (GET) - restricted endpoint to test your access (only for testing)

- `/api/v1/listpods` **RD** (GET) - returns pods from k8s
  - `namespace` (optional) - only list pods in given namespace, when not provided all pods are listed regarding of their namespace
  > deprecated, use `/api/v2/listpods` instead

- `/api/v2/listpods` **T** (GET) - returns the list of pods
  - `namespace` (optional) - only list pods in given namespace, when not provided all pods are listed regarding of their namespace

- `/api/v1/listcontainers` **T** (GET) - returns the list of containers
  - `namespace` (optional) - only list containers in given namespace
  - `pod_name` (optional) - list only containers of given pod

- `/api/v1/listservices` **T** (GET) - returns the list of services
  - `namespace` (optional) - only list services in given namespace

- `/api/v1/listnamespaces` (GET) - lists namespaces of k8s cluster
- `/api/v1/createdeployment` **T** (POST) - create a deployment for the cluster
  - `namespace` - namespace to create deployment in
  - `name` - name of deployment
  - `image` - name of image to be deployed
  - `replicas` - number of replicas (integer 1..32)
  - `cpu_request` - CPU request in k8s format (ex. `100m`)
  - `memory_request` - Memory request in k8s format
  - `cpu_limit` - CPU limit in k8s format (ex. `64Mi`)
  - `memory_limit` - Memory limit in k8s format
- `/api/v1/getpodmetrics` **RTD** (GET) - returns metrics via k8s `metrics server`
  - `namespace` (optional) - only get metrics of pods in given namespace, when not provided all pods are listed regarding of their namespace
  > deprecated, use `/api/v2/getpodmetrics` instead
- `/api/v2/getpodmetrics` **T** (GET) - returns metrics in a given time period
  - `pod_name` (optional) - the name of the pod from which to retrieve metrics data. If this parameter is not provided, metrics for all pods are returned. In case an invalid `pod_name` is provided, the endpoint will still return data, but the `Pods` field in each record will be empty
  - `start_time` (optional) - the start of the metric history time frame in RFC3339 format (ex. `2024-08-24T20:56:12.999Z`) 
  - `end_time` (optional) - the end of the metric history time frame in RFC3339 format, can't be earlier than `start_time`
  > The maximum time range is 4 hours when `pod_name` is not set and 72 otherwise. If both `start_time` and `end_time` are not set, the most recent record will be returned

- `/api/v1/deletepodmetrics` **T** (GET) - deletes metrics from a given time period
  - `start_time` (optional) - the start of the metric history time frame in RFC3339 format (ex. `2024-08-24T20:56:12.999Z`) 
  - `end_time` (optional) - the end of the metric history time frame in RFC3339 format, can't be earlier than `start_time`
  > If only one of `start_time` or `end_time` are provided, the endpoint will delete all metrics newer than `start_time` or older than `end_time`. If none of them are provided, the endpoint will delete all metrics currently recorded


> **R** - the endpoint returns raw data from Kubernetes API without any processing

> **T** - the endpoint needs extensive testing

> **D** - the endpoint is being deprecated

GET requests take Query Parameters while Post requests take Parameters in Body as JSON. Body may be empty it there're no required parameters.

API is not stable and may change at any commit.

## Example

```python
import requests
import sys

API_ADDR = 'http://localhost:5000/api/v1/'

# authenticate
r = requests.post(API_ADDR+'login', json={'user': 'john', 'pass': 'doe'})
if (r.status_code == 200):
    print('login successful')
    token = r.json().get('token')
else:
    sys.exit(2)

# check if default namespace exists
r = requests.get(API_ADDR+'listnamespaces', headers={'Authorization': f'Bearer {token}'}, json={})
if (r.status_code == 200):
    print('got namespaces')
    if 'default' in r.json().get('namespaces'):
        print('default namespace will be used for deployment')
    else:
        print('no default namespace found!')
        sys.exit(1)
else:
    sys.exit(2)

# create deployment

depl_data = {
    'namespace': 'default',
    'name': 'nginx-deploy-with-api',
    'image': 'nginx',
    'replicas': 2,
    'cpu_request': '100m',
    'memory_request': '128Mi',
    'cpu_limit': '200m',
    'memory_limit': '256Mi'
}

r = requests.post(API_ADDR+'createdeployment', headers={'Authorization': f'Bearer {token}'}, json=depl_data)
if (r.status_code == 200):
    print('deployment successful! look at `kubectl get pods --output=wide` on your controller')
```