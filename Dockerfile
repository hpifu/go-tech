FROM centos:centos7
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" >> /etc/timezone
COPY docker/go-tech /var/docker/go-tech
RUN mkdir -p /var/docker/go-tech/log
EXPOSE 6062
WORKDIR /var/docker/go-tech
CMD [ "bin/tech", "-c", "configs/tech.json" ]
