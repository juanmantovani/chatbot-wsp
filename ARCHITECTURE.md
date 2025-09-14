# Arquitectura del Sistema - WhatsApp Chatbot

## Visión General

Este documento describe la arquitectura del sistema de chatbot para WhatsApp Business, implementado en Go siguiendo principios de arquitectura limpia y buenas prácticas de desarrollo.

## Principios de Diseño

### 1. Arquitectura Limpia (Clean Architecture)
- **Separación de responsabilidades**: Cada capa tiene una responsabilidad específica
- **Inversión de dependencias**: Las capas internas no dependen de las externas
- **Independencia de frameworks**: El código de negocio no depende de frameworks específicos

### 2. Patrones de Diseño Implementados
- **Repository Pattern**: Para abstraer el acceso a datos
- **Service Layer**: Para encapsular la lógica de negocio
- **Dependency Injection**: Para facilitar testing y mantenimiento
- **Factory Pattern**: Para la creación de objetos complejos

## Estructura del Proyecto

```
chatbot-wsp/
├── cmd/                          # Punto de entrada de la aplicación
│   └── main.go
├── internal/                     # Código interno de la aplicación
│   ├── domain/                   # Lógica de negocio (Capa de Dominio)
│   │   ├── models/              # Entidades y objetos de valor
│   │   ├── repository/          # Interfaces de repositorio
│   │   ├── service/             # Servicios de dominio
│   │   └── errors/              # Errores personalizados
│   └── infrastructure/          # Infraestructura (Capa de Infraestructura)
│       ├── config/              # Configuración de la aplicación
│       ├── http/                # Capa de presentación HTTP
│       │   ├── handlers/        # Controladores HTTP
│       │   ├── middleware/      # Middlewares HTTP
│       │   └── routes/          # Configuración de rutas
│       └── logger/              # Sistema de logging
├── aws/                         # Configuraciones para AWS
│   ├── cloudformation/          # Templates de CloudFormation
│   └── ecs/                    # Configuraciones de ECS/EKS
├── scripts/                     # Scripts de utilidad
└── tests/                       # Tests de integración
```

## Capas de la Arquitectura

### 1. Capa de Dominio (Domain Layer)
**Ubicación**: `internal/domain/`

**Responsabilidades**:
- Define las reglas de negocio
- Contiene las entidades principales
- Define interfaces para repositorios
- Maneja errores específicos del dominio

**Componentes**:
- `models/`: Entidades como `WhatsAppMessage`, `ChatbotState`, `ChatbotFlow`
- `repository/`: Interfaces para acceso a datos
- `service/`: Lógica de negocio del chatbot
- `errors/`: Errores personalizados del dominio

### 2. Capa de Infraestructura (Infrastructure Layer)
**Ubicación**: `internal/infrastructure/`

**Responsabilidades**:
- Implementa las interfaces definidas en el dominio
- Maneja la configuración externa
- Proporciona servicios de infraestructura

**Componentes**:
- `config/`: Carga y manejo de configuración
- `logger/`: Sistema de logging estructurado
- `http/`: Implementación de la API REST

### 3. Capa de Aplicación (Application Layer)
**Ubicación**: `cmd/main.go`

**Responsabilidades**:
- Orquesta la inicialización de la aplicación
- Configura las dependencias
- Maneja el ciclo de vida de la aplicación

## Flujo de Datos

### 1. Recepción de Mensajes
```
WhatsApp Business API → Webhook → Handler → Service → Repository
```

### 2. Procesamiento de Mensajes
```
Message → State Machine → Flow Selection → Response Generation
```

### 3. Envío de Respuestas
```
Response → Handler → WhatsApp Business API → User
```

## Estados del Chatbot

### 1. Estado Inicial (`welcome`)
- Muestra mensaje de bienvenida
- Presenta 4 opciones (A, B, C, D)
- Espera selección de opción

### 2. Estados de Opción (`option_a`, `option_b`, `option_c`, `option_d`)
- Procesa la opción seleccionada
- Muestra mensaje específico
- Solicita datos relevantes
- Transiciona a estado de recolección

### 3. Estado de Recolección (`collecting_data`)
- Recopila datos del usuario
- Permite seleccionar otra opción
- Mantiene historial de datos

## Configuración y Despliegue

### 1. Configuración Local
- Variables de entorno en `.env`
- Configuración de logging
- Tokens de WhatsApp Business API

### 2. Despliegue en AWS
- **Lambda**: Para arquitectura serverless
- **ECS Fargate**: Para contenedores administrados
- **EKS**: Para orquestación de Kubernetes

### 3. Monitoreo y Logging
- Logs estructurados en JSON
- Integración con CloudWatch
- Métricas de salud y rendimiento

## Patrones de Comunicación

### 1. Webhook de WhatsApp
- **Verificación**: `GET /whatsapp/webhook`
- **Recepción**: `POST /whatsapp/webhook`

### 2. API REST
- **Health Check**: `GET /health`
- **Estadísticas**: `GET /stats`
- **Mensaje de Bienvenida**: `GET /whatsapp/welcome`

## Manejo de Errores

### 1. Estrategia de Errores
- Errores específicos del dominio
- Logging detallado de errores
- Respuestas HTTP apropiadas
- Recuperación graceful

### 2. Tipos de Errores
- `ErrFlowNotFound`: Flujo no encontrado
- `ErrInvalidState`: Estado inválido
- `ErrInvalidOption`: Opción inválida
- `ErrInvalidWebhook`: Webhook inválido

## Testing

### 1. Tests Unitarios
- Tests para servicios de dominio
- Mocks para repositorios
- Cobertura de casos de uso

### 2. Tests de Integración
- Tests de endpoints HTTP
- Tests de flujos completos
- Tests de configuración

## Escalabilidad y Rendimiento

### 1. Escalabilidad Horizontal
- Arquitectura stateless
- Balanceador de carga
- Auto-scaling en AWS

### 2. Optimizaciones
- Pool de conexiones HTTP
- Caching de configuraciones
- Logging asíncrono

## Seguridad

### 1. Autenticación
- Verificación de tokens de WhatsApp
- Validación de webhooks
- Headers de seguridad

### 2. Validación
- Validación de entrada
- Sanitización de datos
- Rate limiting

## Mantenimiento y Evolución

### 1. Extensibilidad
- Fácil adición de nuevos flujos
- Configuración externa de mensajes
- Plugins para nuevas funcionalidades

### 2. Monitoreo
- Métricas de rendimiento
- Alertas automáticas
- Dashboards de monitoreo

## Consideraciones de Producción

### 1. Disponibilidad
- Health checks automáticos
- Restart automático de servicios
- Redundancia en múltiples AZs

### 2. Observabilidad
- Trazabilidad de requests
- Métricas de negocio
- Logs centralizados

### 3. Backup y Recuperación
- Backup de configuraciones
- Estrategia de rollback
- Disaster recovery

Esta arquitectura proporciona una base sólida, escalable y mantenible para el chatbot de WhatsApp Business, siguiendo las mejores prácticas de desarrollo de software.
