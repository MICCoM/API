FROM ubuntu:14.04
MAINTAINER Andreas Wilke <wilke@mcs.anl.gov>
RUN apt-get update && apt-get -y upgrade
# Install basic libs and tools
RUN apt-get -y install ssh sshfs git curl emacs

# Install MongoDB
EXPOSE 27017
RUN sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv EA312927
RUN echo "deb http://repo.mongodb.org/apt/ubuntu trusty/mongodb-org/3.2 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-3.2.list
RUN apt-get update && apt-get install -y mongodb-org
RUN mkdir -p /data/db
# RUN /usr/bin/mongod

# Install GO
EXPOSE 8000 8001
RUN sudo apt-get -y install golang
ENV WORKSPACE=/workspace 
ENV GOPATH=$WORKSPACE
RUN mkdir $WORKSPACE ; cd $WORKSPACE 
RUN go get github.com/MICCoM/API 

