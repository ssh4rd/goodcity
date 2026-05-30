FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
ARG VITE_DADATA_TOKEN=""
ENV VITE_DADATA_TOKEN=$VITE_DADATA_TOKEN
RUN npm run build

FROM node:20-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --production
COPY --from=builder /app/dist ./dist
COPY server.js .
COPY index.html .
ENV NODE_ENV=production
EXPOSE 3000
CMD ["node", "server.js"]
