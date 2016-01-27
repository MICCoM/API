FROM ubuntu:14.04
MAINTAINER Andreas Wilke <wilke@mcs.anl.gov>
RUN apt-get update && apt-get -y upgrade
RUN apt-get -y install ssh sshfs
