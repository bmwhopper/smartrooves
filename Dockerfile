FROM node:8
  
WORKDIR /opt/smartrooves

# install deps
COPY package.json /opt/smartrooves
RUN npm install

# Setup workdir
COPY . /opt/smartrooves

# run
EXPOSE 3000
CMD ["npm", "start"]
