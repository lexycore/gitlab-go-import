# gitlab-go-import


    location / {
        set $destination "http://127.0.0.1:8888";
        if ($arg_go-get = "1") {
            set $destination "http://127.0.0.1:8008";
        }
        client_max_body_size 0;
        gzip off;

        ## https://github.com/gitlabhq/gitlabhq/issues/694
        ## Some requests take more than 30 seconds.
        proxy_read_timeout      300;
        proxy_connect_timeout   300;
        proxy_redirect          off;

        proxy_http_version 1.1;

        proxy_pass              $destination;
        proxy_set_header        X-Forwarded-Proto $scheme;
        proxy_set_header        Host              $http_host;
        proxy_set_header        X-Real-IP         $remote_addr;
        proxy_set_header        X-Forwarded-Ssl   on;
        proxy_set_header        X-Forwarded-For   $proxy_add_x_forwarded_for;
        proxy_set_header        X-Frame-Options   SAMEORIGIN;
    }
