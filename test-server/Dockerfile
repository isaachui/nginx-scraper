FROM nginx
COPY content/ /usr/share/nginx/html

COPY nginx.conf /etc/nginx/nginx.conf

COPY test-server /etc/nginx/sites-enabled/

#allow logs to go to stdout
RUN touch /var/log/nginx/stdout.log
RUN ln -sf /dev/stdout /var/log/nginx/stdout.log
