ARG NGINX_VERSION
FROM nginx:1.17.8

LABEL maintainer="moy"

# replace sources.list

# init
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    wget \
    zip \
    unzip \
    git-core \
    openssl \
    gcc \
    zlib1g-dev \
    libssl-dev \
    libpcre3 \
    libpcre3-dev \
    libc-dev \
    libxslt-dev \
    libxslt1-dev \
    libgd-dev \
    libgeoip-dev \
    libpcre++-dev \
    libssl-dev \
    libperl-dev \
    perl \
    build-essential \
    ca-certificates \
    procps \
    tmux \
    keychain \
    net-tools \
    curl \
    make \
    vim \
    locales \
    expect \
    cron \
    && rm -rf /var/lib/apt/lists/* \
    && mkdir -p /home/downloads/

# nginx-release
RUN cd /home/downloads/ \
    && wget --no-check-certificate https://github.com/nginx/nginx/archive/release-1.17.8.zip -O nginx-src.zip \
    && unzip nginx-src.zip \
    && rm -f nginx-src.zip \
    # nginx-upload-module-2.3.0
    && git clone https://github.com/fdintino/nginx-upload-module.git \
    # nginx-upload-progress-module-0.9.2
    && git clone https://github.com/masterzen/nginx-upload-progress-module.git \
    && cd /home/downloads/nginx-release-1.17.8 \
    && ./auto/configure --prefix=/etc/nginx --sbin-path=/usr/sbin/nginx \
    --modules-path=/usr/lib/nginx/modules --conf-path=/etc/nginx/nginx.conf --error-log-path=/var/log/nginx/error.log \
    --http-log-path=/var/log/nginx/access.log --pid-path=/var/run/nginx.pid --lock-path=/var/run/nginx.lock \
    --http-client-body-temp-path=/var/cache/nginx/client_temp --http-proxy-temp-path=/var/cache/nginx/proxy_temp \
    --http-fastcgi-temp-path=/var/cache/nginx/fastcgi_temp --http-uwsgi-temp-path=/var/cache/nginx/uwsgi_temp \
    --http-scgi-temp-path=/var/cache/nginx/scgi_temp --user=nginx --group=nginx --with-compat --with-file-aio --with-threads \
    --with-http_addition_module --with-http_auth_request_module --with-http_dav_module --with-http_flv_module --with-http_gunzip_module \
    --with-http_gzip_static_module --with-http_mp4_module --with-http_random_index_module --with-http_realip_module --with-http_secure_link_module \
    --with-http_slice_module --with-http_ssl_module --with-http_stub_status_module --with-http_sub_module --with-http_v2_module \
    --with-mail --with-mail_ssl_module --with-stream --with-stream_realip_module --with-stream_ssl_module --with-stream_ssl_preread_module \
    --with-http_xslt_module=dynamic \
    --with-http_image_filter_module=dynamic \
    --with-http_geoip_module=dynamic \
    --with-http_perl_module=dynamic \
    # --with-ipv6 \
    --add-module=/home/downloads/nginx-upload-progress-module \
    --add-module=/home/downloads/nginx-upload-module \
    && make -j$(nproc) \
    && make install

RUN echo 'en_US.UTF-8 UTF-8' >> /etc/locale.gen && locale-gen \
    && mkdir -p /usr/local/nginx/logs/ && mkdir -p /works \
    && service cron stop \
    && service cron start

ENV LANG=en_US.UTF-8 \
    LC_ALL=en_US.UTF-8

STOPSIGNAL SIGTERM

CMD ["nginx", "-g", "daemon off;"]
