proxy_cache_path /var/cache/nginx/pkgdex_cache
  levels=1:2
  keys_zone=pkgdex_cache:10m
  max_size=10g
  inactive=60m
  use_temp_path=off;

limit_req_zone $binary_remote_addr zone=pkgdex_api:10m rate=10r/s;
limit_req_zone $binary_remote_addr zone=pkgdex_general:10m rate=30r/s;

upstream pkgdex_backend {
  server 127.0.0.1:1997;
  keepalive 32;
}

server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name pkg.example.com;
    root /srv/pkg.example.com;

    log_format cache_status '$remote_addr - $remote_user [$time_local] '
                           '"$request" $status $body_bytes_sent '
                           '"$http_referer" "$http_user_agent" '
                           'cache:$upstream_cache_status';

    access_log /var/log/nginx/pkgdex.access.log cache_status buffer=512k flush=1m;
    error_log /var/log/nginx/pkgdex.error.log warn;

    ssl_certificate /etc/nginx/ssl/pkg.example.com.crt;
    ssl_certificate_key /etc/nginx/ssl/pkg.example.com.key;

    location /assets/ {
        expires max;

        limit_req zone=pkgdex_general burst=60 nodelay;

        proxy_cache pkgdex_cache;
        proxy_cache_key $request_uri;
        proxy_cache_valid 200 1d;
        proxy_cache_use_stale error timeout updating http_500 http_502 http_503 http_504;
        proxy_cache_background_update on;
        proxy_cache_lock on;
        add_header Cache-Status $upstream_cache_status;

        proxy_pass https://pkgdex_backend;

        proxy_set_header Host $http_host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        proxy_buffering on;
        proxy_buffers 8 16k;
        proxy_buffer_size 16k;

        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;

        proxy_redirect off;
        proxy_http_version 1.1;
    }

    location = /feed.xml {
        expires 1h;

        limit_req zone=pkgdex_general burst=10 nodelay;

        proxy_cache pkgdex_cache;
        proxy_cache_key $request_uri;
        proxy_cache_valid 200 1d;
        proxy_cache_use_stale error timeout updating http_500 http_502 http_503 http_504;
        proxy_cache_background_update on;
        proxy_cache_lock on;
        add_header Cache-Status $upstream_cache_status;

        proxy_pass https://pkgdex_backend;

        proxy_set_header Host $http_host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        proxy_buffering on;
        proxy_buffers 8 16k;
        proxy_buffer_size 16k;

        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;

        proxy_redirect off;
        proxy_http_version 1.1;
    }

    location = /sitemap.xml {
        expires 1d;

        limit_req zone=pkgdex_general burst=10 nodelay;

        proxy_cache pkgdex_cache;
        proxy_cache_key $request_uri;
        proxy_cache_valid 200 1d;
        proxy_cache_use_stale error timeout updating http_500 http_502 http_503 http_504;
        proxy_cache_background_update on;
        proxy_cache_lock on;
        add_header Cache-Status $upstream_cache_status;

        proxy_pass http://pkgdex_backend;

        proxy_set_header Host $http_host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        proxy_buffering on;
        proxy_buffers 8 16k;
        proxy_buffer_size 16k;

        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;

        proxy_redirect off;
        proxy_http_version 1.1;
    }

    location /meta/ {
        limit_req zone=pkgdex_api burst=5 nodelay;

        proxy_cache off;

        proxy_pass https://pkgdex_backend;

        proxy_set_header Host $http_host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        proxy_buffering on;
        proxy_buffers 8 16k;
        proxy_buffer_size 16k;

        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;

        proxy_redirect off;
        proxy_http_version 1.1;
    }

    location / {
        limit_req zone=pkgdex_general burst=30 nodelay;

        set $cache_bypass 0;
        if ($arg_go-get = '1') {
          set $cache_bypass 1;
        }

        proxy_cache pkgdex_cache;
        proxy_no_cache $cache_bypass;
        proxy_cache_bypass $cache_bypass;
        proxy_cache_key $request_uri;
        proxy_cache_valid 200 30m;
        proxy_cache_use_stale error timeout updating http_500 http_502 http_503 http_504;
        proxy_cache_background_update on;
        proxy_cache_lock on;
        add_header Cache-Status $upstream_cache_status;

        proxy_pass https://pkgdex_backend;

        proxy_set_header Host $http_host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        proxy_buffering on;
        proxy_buffers 8 16k;
        proxy_buffer_size 16k;

        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;

        proxy_redirect off;
        proxy_http_version 1.1;
    }
}

server {
  listen 80;
  listen [::]:80;
  server_name pkg.example.com;

  location / {
    return 301 https://$server_name$request_uri;
  }
}
