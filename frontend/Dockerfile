# Etapa 1: Construcción de la aplicación con Vite
FROM node:22 as build

# Establece el directorio de trabajo en el contenedor
WORKDIR /home/app

# Copia el package.json y package-lock.json
COPY package*.json ./

# Instala las dependencias del proyecto
RUN npm install

# Copia el resto de los archivos del proyecto
COPY . .

# Define el argumento de compilación
ARG API_URL

# Establece la variable de entorno
ENV API_URL=${API_URL}

# Genera el archivo .env
RUN node generate-env.js

# Construye el proyecto
RUN npm run build

# Etapa 2: Servir la aplicación con Serve
RUN npm install -g serve

# Utiliza el puerto 5173
EXPOSE 5173

# Inicia la aplicación
CMD ["serve", "-s", "dist", "-l", "5173"]
