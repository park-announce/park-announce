# Start from the official Node.js base image
FROM node:14-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the package.json and package-lock.json files
COPY package.json ./

# Install the dependencies
#RUN npm ci --only=production
RUN npm i
# Copy the source code into the container
COPY . .

# Build the React app
RUN npm run build

# Expose the desired port
EXPOSE 5173

# Set the command to start the app
CMD ["npm", "start"]