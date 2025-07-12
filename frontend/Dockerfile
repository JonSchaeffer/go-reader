FROM node:18-alpine

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Create necessary directories and set proper ownership
RUN mkdir -p .svelte-kit tmp /tmp/vite-cache && \
  chown -R node:node /app /tmp/vite-cache

# Switch to node user
USER node

# Expose port
EXPOSE 3000

# Start development server
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0", "--port", "3000"]
