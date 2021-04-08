gitlab-go-import
================

# why

When we use project groups in GitLab, golang is unable to find some imports
until we add `.git` suffix to module path. This approach isn't neat at all.

This small app handles requests with `go-get=1` query argument and resolves
such cases enabling clean usage if import paths.

# how

For GitLab installations running behind nginx proxy, server configuration
could be changed to something like this:


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

This means requests with `go-get=1` query argument will be handled by this
application with is serving on `127.0.0.1:8008`.

The best way to install this application is:

    cd /opt
    git clone https://github.com/lexycore/gitlab-go-import.git
    cd ./gitlab-go-import

    # edit docker-compose.yml

    docker-compose up --build


Of course, we need to configure this app. This can be done via `docker-compose.yml`.

 - `GO_IMPORT_GITLAB_URL` [env var] - url of your GitLab server
 - `GO_IMPORT_GITLAB_TOKEN` [env var] - private access token created by GitLab admin user
 - `ports: ["8008:8008"]` - configure port on which nginx will be redirecting requests to this app
