FROM centos:8
WORKDIR /app
COPY wee-server /app/
COPY configuration.json /app/
RUN rm -f /etc/localtime \
&& ln -sv /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
&& echo "Asia/Shanghai" > /etc/timezone

EXPOSE 8080/tcp
CMD /app/wee-server -c /app/configuration-prod.json