# 🚀 Go Fiber + PostgreSQL on Kubernetes

This project demonstrates a complete deployment of a **Go Fiber REST API** connected to a **PostgreSQL** database on **Kubernetes**.  
The setup uses `ConfigMap`, `Secret`, `Deployment`, `StatefulSet`, and `Service`, all managed via **Kustomize**.

---

## 🧱 Project Structure
```` 
kubernetes/
├── go-fiber/
│ ├── configmap.yaml
│ ├── deployment.yaml
│ ├── kustomization.yaml
│ ├── secret.yaml
│ └── service.yaml
└── postgre/
├── configmap.yaml
├── kustomization.yaml
├── secret.yaml
├── service.yaml
└── statefulset.yaml
```` 

## 🧩 Architecture Overview
````
+-------------------+ +-------------------+
| Go Fiber App | <----> | PostgreSQL DB |
| (go-fiber-deploy) | | (statefulset) |
| Port: 3000 | | Port: 5432 |
+-------------------+ +-------------------+
↑ ↑
ConfigMap & Secret ConfigMap & Secret
↓ ↓
Environment Vars Environment Vars
````

- The **Go Fiber** service (`go-fiber-http`) communicates internally with the **PostgreSQL** service (`postgres`) inside the `go-fiber` namespace.
- Configuration and credentials are injected via **ConfigMaps** and **Secrets**.

---

## ⚙️ Prerequisites

Before deploying, ensure you have:
- [Docker](https://www.docker.com/)
- [Minikube](https://minikube.sigs.k8s.io/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Kustomize](https://kubectl.docs.kubernetes.io/installation/kustomize/)

---

## 🗂️ Kubernetes Resources

### 🐘 PostgreSQL (Database)
Located in: `kubernetes/postgre/`

- **ConfigMap** → database host & port  
- **Secret** → database username, password, and DB name  
- **StatefulSet** → PostgreSQL pod with persistent volume  
- **Service** → internal ClusterIP service  
- **Kustomization** → combines all manifests for easy apply/delete  

Example deployment command:
````
kubectl apply -k kubernetes/postgre/
```` 
## ⚡ Go Fiber (Application)

Located in: ```` kubernetes/go-fiber/````

**Deployment** → runs Go Fiber REST API pods (replicas: 3)

**ConfigMap** → shared app configuration

**Secret** → reused PostgreSQL credentials

**Service** → ClusterIP exposing API internally (port 3000)

**Kustomization** → combines all manifests

Example deployment command:
`````
kubectl apply -k kubernetes/go-fiber/
```````
## 🚀 Deployment Guide

**1. Start Minikube** 
````
kubectl apply -k kubernetes/go-fiber/
````
**2. Create Namespace** 
````
kubectl create namespace go-fiber
````
**3. Deploy PostgreSQL** 
````
kubectl apply -k kubernetes/postgre/
````
**4. Wait for PostgreSQL Pod to Run**
````
kubectl get pods -n go-fiber
````
**5. Deploy Go Fiber App**
````
kubectl apply -k kubernetes/go-fiber/
````
**6. Verify All Resources**
````
kubectl get all -n go-fiber
````

## 🔍 Testing the Application
**1. Port Forward the Go Fiber Service**
````
kubectl port-forward svc/go-fiber-http 3000:3000 -n go-fiber
````
**2. Access API Endpoint**
Visit:
````
http://127.0.0.1:3000/api/categories
````
**3. Expected Output**
````
{
  "data": [],
  "message": "Category list retrieved successfully"
}
````

## 🧠 Environment Variables Overview
| Variable      | Description             | Source    |
| ------------- | ----------------------- | --------- |
| `DB_HOST`     | PostgreSQL service name | ConfigMap |
| `DB_PORT`     | PostgreSQL port (5432)  | ConfigMap |
| `DB_USER`     | Database username       | Secret    |
| `DB_PASSWORD` | Database password       | Secret    |
| `DB_NAME`     | Database name           | Secret    |

## 🧩 Example Docker Image
Your Go Fiber app uses the following image (built from your app source):
````
image: repodocker/go-fiber-example:latest
````
You can build and push updates using:
````
docker build -t repodocker/go-fiber-example:latest .
docker push repodocker/go-fiber-example:latest
````
## 🧹 Cleanup
To delete all deployed resources:
````
kubectl delete -k kubernetes/go-fiber/
kubectl delete -k kubernetes/postgre/
````
Or delete the namespace entirely:
````
kubectl delete namespace go-fiber
````
## 🧾 Notes
• Use StatefulSet for PostgreSQL to ensure persistent data storage.

• Update Secret values with your own credentials for production use.

• The setup uses ClusterIP networking; for external access, use an Ingress or LoadBalancer type.

• Ensure both app and database share the same namespace: go-fiber.

