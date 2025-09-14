# WhatsApp Business Chatbot API - BabyHome

Una aplicación en Go que funciona como chatbot médico para WhatsApp Business de la Dra. Carla Narváez - BabyHome, implementando arquitectura limpia y buenas prácticas de desarrollo.

## Características

- 🤖 Chatbot médico con 4 opciones especializadas (A, B, C, D)
- 👩‍⚕️ Servicios de la Dra. Carla Narváez - BabyHome
- 📱 Integración con WhatsApp Business API
- 🏗️ Arquitectura limpia y escalable
- 🐳 Containerización con Docker
- ☁️ Despliegue en AWS (Lambda, ECS, EKS)
- 📊 Logging estructurado
- 🔒 Manejo seguro de tokens
- 🧪 Cobertura de tests

## Arquitectura

```
cmd/
├── main.go                    # Punto de entrada de la aplicación

internal/
├── domain/                    # Lógica de negocio
│   ├── models/               # Modelos de datos
│   ├── repository/           # Interfaces de repositorio
│   ├── service/              # Servicios de negocio
│   └── errors/               # Errores personalizados
├── infrastructure/           # Infraestructura
│   ├── config/              # Configuración
│   ├── http/                # Handlers HTTP
│   │   ├── handlers/        # Controladores
│   │   ├── middleware/      # Middlewares
│   │   └── routes/          # Rutas
│   └── logger/              # Sistema de logging
└── mocks/                   # Mocks para testing

aws/                         # Configuraciones de AWS
├── cloudformation/          # Templates de CloudFormation
└── ecs/                    # Configuraciones de ECS/EKS

scripts/                     # Scripts de despliegue
```

## Flujo del Chatbot

1. **Saludo inicial**: El usuario recibe un mensaje de bienvenida
2. **Opciones**: Se presentan 4 opciones (A, B, C, D)
3. **Procesamiento**: Cada opción lleva a un flujo específico
4. **Recolección de datos**: Se solicitan datos relevantes según la opción
5. **Continuación**: El usuario puede seleccionar otra opción

### Opciones disponibles:
- **A**: Realizar consulta médica telefónica ($15.000 ARS)
- **B**: Enviar estudios para lectura ($15.000 ARS)
- **C**: Solicitar turno en consultorio
- **D**: Consulta sobre BabyHome

## Instalación y Configuración

### Prerrequisitos

- Go 1.21 o superior
- Docker (opcional)
- AWS CLI (para despliegue)

### Configuración local

1. **Clonar el repositorio**:
```bash
git clone <repository-url>
cd chatbot-wsp
```

2. **Instalar dependencias**:
```bash
make deps
```

3. **Configurar variables de entorno**:
```bash
cp env.example .env
# Editar .env con tus valores
```

4. **Ejecutar la aplicación**:
```bash
make run
```

### Variables de entorno

```env
# Server Configuration
PORT=8080
HOST=0.0.0.0

# WhatsApp Business API
WHATSAPP_VERIFY_TOKEN=your_verify_token_here
WHATSAPP_ACCESS_TOKEN=your_access_token_here

# AWS Configuration
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key

# Logging
LOG_LEVEL=info
```

## Uso con Docker

### Construir imagen
```bash
make docker-build
```

### Ejecutar contenedor
```bash
make docker-run
```

### Con docker-compose
```bash
docker-compose up -d
```

## Despliegue en AWS

### Opción 1: Lambda (Serverless)

```bash
# Configurar variables de entorno
export WHATSAPP_VERIFY_TOKEN="your_token"
export WHATSAPP_ACCESS_TOKEN="your_token"
export ENVIRONMENT="dev"

# Desplegar
make deploy
```

### Opción 2: ECS Fargate

1. **Construir y subir imagen a ECR**:
```bash
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com
docker build -t chatbot-wsp .
docker tag chatbot-wsp:latest <account-id>.dkr.ecr.us-east-1.amazonaws.com/chatbot-wsp:latest
docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/chatbot-wsp:latest
```

2. **Crear cluster ECS y desplegar**:
```bash
aws ecs create-cluster --cluster-name chatbot-wsp-cluster
# Actualizar task-definition.json con tu account-id
aws ecs register-task-definition --cli-input-json file://aws/ecs/task-definition.json
```

### Opción 3: EKS (Kubernetes)

```bash
# Aplicar configuraciones de Kubernetes
kubectl apply -f aws/ecs/service.yaml
```

## Endpoints de la API

### Webhook de WhatsApp
- `GET /whatsapp/webhook` - Verificación del webhook
- `POST /whatsapp/webhook` - Recibir mensajes de WhatsApp

### Endpoints de utilidad
- `GET /health` - Health check
- `GET /stats` - Estadísticas del servicio
- `GET /whatsapp/welcome` - Mensaje de bienvenida

## Configuración del Webhook de WhatsApp

1. **Configurar webhook en Meta for Developers**:
   - URL: `https://tu-dominio.com/whatsapp/webhook`
   - Verify Token: El valor de `WHATSAPP_VERIFY_TOKEN`
   - Webhook Fields: `messages`

2. **Verificar configuración**:
```bash
curl "https://tu-dominio.com/whatsapp/webhook?hub.mode=subscribe&hub.verify_token=TU_TOKEN&hub.challenge=CHALLENGE"
```

## Testing

```bash
# Ejecutar todos los tests
make test

# Ejecutar tests con cobertura
go test -v -cover ./...

# Ejecutar tests específicos
go test -v ./internal/domain/service/
```

## Desarrollo

### Estructura de commits
```
feat: add new feature
fix: fix bug
docs: update documentation
style: formatting changes
refactor: code refactoring
test: add tests
chore: maintenance tasks
```

### Flujo de trabajo
1. Crear rama feature: `git checkout -b feature/nueva-funcionalidad`
2. Hacer cambios y commits
3. Ejecutar tests: `make test`
4. Crear pull request

## Monitoreo y Logging

La aplicación incluye logging estructurado con los siguientes niveles:
- `DEBUG`: Información detallada para debugging
- `INFO`: Información general de la aplicación
- `WARN`: Advertencias que no detienen la ejecución
- `ERROR`: Errores que requieren atención

### Logs en AWS
- **CloudWatch Logs**: Para aplicaciones Lambda y ECS
- **CloudTrail**: Para auditoría de API calls
- **X-Ray**: Para tracing distribuido (opcional)

## Contribución

1. Fork el proyecto
2. Crear una rama feature (`git checkout -b feature/AmazingFeature`)
3. Commit los cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abrir un Pull Request

## Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo `LICENSE` para más detalles.

## Soporte

Para soporte, por favor abrir un issue en el repositorio o contactar al equipo de desarrollo.
