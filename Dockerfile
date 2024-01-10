FROM nginx:alpine AS base
RUN rm /etc/nginx/conf.d/default.conf

FROM base as development
COPY nginx-dev.conf /etc/nginx/conf.d/default.conf

FROM base
COPY nginx.conf /etc/nginx/conf.d/default.conf
