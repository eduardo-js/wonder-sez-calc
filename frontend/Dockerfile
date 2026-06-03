# Frontend dev-server image for wonderlic-calc.
# Runs the Vite dev server (with HMR) bound to 0.0.0.0 so the host can reach it.
FROM node:24-alpine
WORKDIR /app

# Install deps first for layer caching. Source is bind-mounted at runtime (see docker-compose.yml).
COPY package.json package-lock.json ./
RUN npm install

COPY . .

EXPOSE 5173
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0", "--port", "5173"]
