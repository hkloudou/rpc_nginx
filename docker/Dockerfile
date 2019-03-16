FROM alpine:3.9
COPY /bin/rpc-nginx /bin/rpc-nginx
#RUN apk add --no-cache \
#		ca-certificates
#RUN apk del .build-deps

ENV GRPC_PORT :5000
ENV VIRTUAL_PROTO grpc
EXPOSE 5000
CMD ["rpc-nginx"]