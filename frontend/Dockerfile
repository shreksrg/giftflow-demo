# Build Stage
FROM public.ecr.aws/docker/library/node:18-alpine as build-stage

WORKDIR /app

# Copy package files first to leverage Docker cache
COPY package*.json ./

# Use BuildKit cache mount to cache npm packages and speed up install
# Use npm ci for faster, reliable installs from lockfile
RUN --mount=type=cache,target=/root/.npm \
    npm ci

COPY . .
RUN npm run build

# Production Stage
FROM public.ecr.aws/nginx/nginx:stable-alpine as production-stage

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build-stage /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
