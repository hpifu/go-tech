FROM centos:centos7
RUN timedatectl set-timezone Asia/Shanghai
COPY docker/go-tech /var/docker/go-tech
RUN mkdir -p /var/docker/go-tech/log
EXPOSE 6062
WORKDIR /var/docker/go-tech
CMD [ "bin/tech", "-c", "configs/tech.json" ]
