# Use the official Node image as the base image
FROM node:20 as build

# Set the working directory in the container
WORKDIR /app

# Copy the package.json and package-lock.json files to the working directory
COPY package.json package-lock.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application files to the working directory
COPY . ./

# Build the React app
RUN npm run build

# Use a lightweight Nginx image as the base image for serving the React app
FROM nginx:alpine

# Copy the built React app from the previous stage to the Nginx server's public directory
COPY --from=build /app/build /usr/share/nginx/html

# Expose port 80
EXPOSE 80

# Start Nginx when the container starts
CMD ["nginx", "-g", "daemon off;"]
