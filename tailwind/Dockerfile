FROM node:latest
WORKDIR /tailwind

# Install dependencies
RUN npm init -y && \
       npm install tailwindcss postcss-cli autoprefixer && \
         npx tailwindcss init;

CMD npx tailwindcss -c /src/tailwind.config.js -i /src/styles.css -o /dst/styles.css --watch --poll
