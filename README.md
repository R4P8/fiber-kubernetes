## 🚀 Go Fiber + PostgreSQL + Observability (Kubernetes Ready)

Proyek ini merupakan implementasi lengkap dari **REST API menggunakan Go Fiber**  yang terhubung ke **PostgreSQL**, serta dilengkapi dengan **OpenTelemetry, Prometheus, Jaeger**, untuk observability.
Didesain untuk berjalan baik secara **lokal menggunakan Docker Compose**, maupun di Kubernetes.

## 🧱 Project Structure 
````
.
├── config/               # Konfigurasi database dan environment
├── controllers/          # Handler API
├── entities/             # Struktur data dan model
├── repository/           # Layer akses database
├── routes/               # Route / endpoint API
├── .env.example          # Contoh environment variable
├── Dockerfile            # Build Go Fiber app image
├── docker-compose.yml    # Local setup (Postgres + App + Monitoring)
├── go.mod / go.sum       # Dependensi Go module
├── main.go               # Entry point aplikasi
├── otel-config.yaml      # Konfigurasi OpenTelemetry Collector
└── prometheus.yaml       # Konfigurasi Prometheus
````
## ⚡ Features
**✅ Go Fiber Framework** — ringan, cepat, dan efisien
**✅ PostgreSQL** — database relational utama
**✅ OpenTelemetry** Collector — mengumpulkan metrics dan traces
**✅ Prometheus** — metrics monitoring
**✅ Jaeger** — distributed tracing visualizer
**✅ Grafana(Tambahkan)** — observability dashboard
**✅ Docker & Kubernetes** Ready — bisa dijalankan di kedua environment

## 🧩 Architecture Overview
````
Go Fiber App ──> OpenTelemetry Collector ──> Jaeger (Traces)
                     │
                     └────> Prometheus ──> Grafana (Metrics)
````

## ⚙️ Setup (Local with Docker Compose)
**1 Salin file** ````.env.example````
````
cp .env.example .env
````
**2 Build dan Jalankan** 
````
docker-compose up --build
````
**3 Cek container**
````
docker ps
````
**4 Akses layanan**
| Service      | URL                                              |
| ------------ | ------------------------------------------------ |
| Go Fiber API | [http://localhost:3000](http://localhost:3000)   |
| Prometheus   | [http://localhost:9090](http://localhost:9090)   |
| Grafana      | [http://localhost:3001](http://localhost:3001)   |
| Jaeger UI    | [http://localhost:16686](http://localhost:16686) |

## 🧩 Observability Stack
**🟢 OpenTelemetry ````(otel-config.yaml)````**
Mengumpulkan trace dan metrics dari aplikasi Go Fiber lalu mengirim ke:
* Jaeger (trace): ````jaeger:4317````
* Prometheus (metrics): ````0.0.0.0:9464````

**🔵 Prometheus ````(prometheus.yaml)````**
Mengambil metrics dari:

* OpenTelemetry Collector ````(otel-collector:9464)````
* Node Exporter ````(node-exporter:9100)````
* Prometheus sendiri ````(localhost:9090)````

**🟣 Jaeger**
Menampilkan trace aplikasi Go Fiber yang dikirim melalui OTLP port ````4317````.

**🟠 Grafana(Tambahkan Manual)**
Datasource otomatis:
````
- name: Jaeger
  type: jaeger
  url: http://jaeger:16686
  isDefault: true
- name: Prometheus
  type: prometheus
  url: http://prometheus:9090
  isDefault: false
`````
Gunakan Grafana untuk memantau:

* Request latency
* Application traces 
* Resource metrics (CPU, memory, dsb.)

## 🧱 API Structure

| Method   | Endpoint              | Description           |
| -------- | --------------------- | --------------------- |
| `GET`    | `/api/categories`     | Get all categories    |
| `POST`   | `/api/categories`     | Create a new category |
| `PUT`    | `/api/categories/:id` | Update category       |
| `DELETE` | `/api/categories/:id` | Delete category       |

## 🐳 Docker Commands

**Build image**

``` 
docker build -t fiber-kubernetes:latest . 
```

**Run container**
```
docker run -p 3000:3000 fiber-kubernetes:latest
````

## ☸️ Kubernetes Deployment
Tersedia di folder:
```
kubernetes/
├── go-fiber/
├── postgre/
├── otel-collector/
├── prometheus/
├── grafana/
└── jaeger/
```

Langkah deploy:
```
kubectl create namespace go-fiber
kubectl apply -k kubernetes/postgre/
kubectl apply -k kubernetes/go-fiber/
kubectl apply -k kubernetes/otel-collector/
kubectl apply -k kubernetes/prometheus/
kubectl apply -k kubernetes/jaeger/
kubectl apply -k kubernetes/grafana/
````

## 🧪 Test Endpoint 
Port-forward service:
````
kubectl port-forward svc/go-fiber-http 3000:3000 -n go-fiber
`````
Expected response:
````
{
  "data": [],
  "message": "Category list retrieved successfully"
}
``````

## 🧹 Cleanup
Hapus semua resource:
`````
kubectl delete namespace go-fiber
`````
## 🧾 Notes

* Gunakan **StatefulSet** untuk PostgreSQL agar data persisten.
* Semua observability tools berjalan di namespace yang sama. 
* Untuk produksi, ubah kredensial di ``Secret`` dan konfigurasi environment di ``.env.``

