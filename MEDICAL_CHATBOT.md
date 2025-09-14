# Chatbot Médico BabyHome - Dra. Carla Narváez

## Descripción del Servicio

Este chatbot está diseñado específicamente para el servicio médico de la Dra. Carla Narváez - BabyHome, especializado en pediatría y atención neonatal.

## Servicios Ofrecidos

### 🤖 Chatbot BabyHome – Dra. Carla Narváez

**Mensaje de bienvenida:**
```
🤖 Chatbot BabyHome – Dra. Carla Narváez
👋 ¡Hola! Gracias por comunicarte.
Por favor, seleccioná una opción escribiendo la letra correspondiente:
A. Realizar consulta médica telefónica
B. Enviar estudios para lectura
C. Solicitar turno en consultorio
D. Consulta sobre BabyHome
(Si es una urgencia, por favor acudí a una guardia)
```

## Opciones del Chatbot

### A. Realizar consulta médica telefónica

**Costo:** $15.000 ARS (no cubierta por obra social)

**Información requerida:**
1️⃣ Nombre y edad del paciente
2️⃣ Motivo de la consulta
3️⃣ Comprobante de pago (Alias: Narvaez.Carla.B)

**Información importante:**
- https://appar.com.ar/consulta-pediatrica-online/
- 📌 Una vez completados estos pasos, la Dra. se pondrá en contacto

### B. Enviar estudios para lectura

**Costo:** $15.000 ARS

**Información requerida:**
1️⃣ Fotos claras o PDF de los estudios
2️⃣ Síntomas actuales y fecha de realización
3️⃣ Tu duda o pregunta principal
4️⃣ Comprobante de pago (Alias: Narvaez.Carla.B)

**Información importante:**
- https://appar.com.ar/consulta-pediatrica-online/
- 📌 Una vez completados estos pasos, la Dra. se pondrá en contacto

### C. Solicitar turno en consultorio

**Contactos para turnos:**
- Centro Médico Cervantes (WhatsApp: 343-4066281)
- Consultorios OSPEP (WhatsApp: 343-5138637)

### D. Consulta sobre BabyHome

**Servicios de BabyHome:**
💜 ¡Qué alegría que te interese BabyHome!

**Ofrecemos:**
✅ Consulta prenatal
✅ Recepción neonatal personalizada (COPAP y primera hora siempre que mamá y bebé estén clínicamente bien)
✅ Controles en domicilio

**Para orientarte, contanos:**
1️⃣ Semana de embarazo / FPP
2️⃣ Maternidad y obstetra
3️⃣ Si desean priorizar COPAP/primera hora
4️⃣ Si quieren coordinar una consulta prenatal

## Flujo de Conversación

### 1. Estado Inicial
- Usuario recibe mensaje de bienvenida
- Se presentan las 4 opciones (A, B, C, D)
- Usuario selecciona una opción

### 2. Procesamiento de Opción
- Se muestra información específica de la opción seleccionada
- Se solicitan datos relevantes según el servicio
- Se proporciona información de contacto y pagos

### 3. Recolección de Datos
- El chatbot recopila la información proporcionada
- Mantiene un registro de los datos del usuario
- Permite seleccionar otra opción si es necesario

### 4. Continuación
- Usuario puede seleccionar otra opción
- Se mantiene el historial de la conversación
- Se puede volver al menú principal

## Características Técnicas

### Emojis y Formato
- Uso de emojis médicos y profesionales
- Formato optimizado para WhatsApp
- Mensajes claros y concisos
- Información estructurada con numeración

### Datos Recopilados
- **Consulta médica**: Datos del paciente y motivo
- **Lectura de estudios**: Información clínica y estudios
- **Turnos**: Datos de contacto
- **BabyHome**: Información prenatal y neonatal

### Integración WhatsApp
- Webhook para recepción de mensajes
- Respuestas automáticas
- Manejo de estados de conversación
- Logging de interacciones

## Configuración para Producción

### Variables de Entorno
```env
# WhatsApp Business API
WHATSAPP_VERIFY_TOKEN=your_verify_token_here
WHATSAPP_ACCESS_TOKEN=your_access_token_here

# Configuración del servidor
PORT=8080
HOST=0.0.0.0
LOG_LEVEL=info
```

### Webhook de WhatsApp
- **URL**: `https://tu-dominio.com/whatsapp/webhook`
- **Verificación**: `GET /whatsapp/webhook`
- **Recepción**: `POST /whatsapp/webhook`

## Monitoreo y Logs

### Logs Estructurados
- Interacciones de usuarios
- Selección de opciones
- Datos recopilados
- Errores y excepciones

### Métricas Importantes
- Número de consultas por opción
- Tiempo de respuesta
- Datos recopilados por servicio
- Errores de procesamiento

## Consideraciones Médicas

### Urgencias
- El chatbot incluye advertencia sobre urgencias
- Redirige a guardias médicas cuando es necesario
- No reemplaza la atención médica de emergencia

### Privacidad
- Los datos médicos se manejan con confidencialidad
- Cumplimiento con normativas de privacidad médica
- Almacenamiento seguro de información sensible

### Profesionalismo
- Mensajes médicos apropiados
- Información clara sobre costos
- Contactos profesionales verificados

## Mantenimiento

### Actualizaciones de Contenido
- Precios de servicios
- Información de contacto
- Nuevos servicios
- Horarios de atención

### Monitoreo Continuo
- Funcionamiento del webhook
- Respuestas del chatbot
- Calidad de las interacciones
- Rendimiento del sistema

Este chatbot está diseñado para proporcionar una experiencia profesional y eficiente para los pacientes de la Dra. Carla Narváez, facilitando el acceso a sus servicios médicos especializados.
