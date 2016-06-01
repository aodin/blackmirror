server {
    listen 80;

    server_name _;

    location / {
        proxy_pass          http://127.0.0.1:{{ port }};
        proxy_redirect      off;
        proxy_buffering     off;
        proxy_set_header    Host            $host;
        proxy_set_header    X-Real-IP       $remote_addr;
        proxy_set_header    X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
