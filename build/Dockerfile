FROM alpine:3.15
RUN set -eux &&sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add nginx
COPY nginx.conf /etc/nginx/nginx.conf
ADD ./package/web /web
ADD ./package/service/openaiproxy /services/openaiproxy/openaiproxy
copy entrypoint.sh /entrypoint.sh
ENTRYPOINT ["sh","entrypoint.sh"]