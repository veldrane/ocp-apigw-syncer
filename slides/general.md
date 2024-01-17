---
title: Syncer
separator: <!--s-->
verticalSeparator: <!--v-->
revealOptions:
transition: 'none'
---

<!-- .slide: data-background="images/the-problem-background-2.png" data-background-size="1920px" -->

<div id=left3-white-bg>

The problem

* Nginx Plus provides distributed key/val store but without consistency
* there is no guarantee that token is stored in all members (pods) (used like token storage)
* Service object in k8s/openshift distribute requests in random way

<BR>
<BR>

<em>
Fast following request send by frontend can be denied by apigw, even it has a valid token!
</em>

</div>

<!--s-->

<!-- .slide: data-background="images/login-current7.jpg" data-background-size="1920px" -->

<div id=left2-small>

#### Current solution

</div>

<!--s-->

<!-- .slide: data-background="images/default-text-background-black-text-2.png" data-background-size="1920px" -->

<div id=left3-small>

#### Current solution

* we expect that delay between store and send token is enough for synchronization with other pod 
* easy solution, but still without any warranty
* dark side of the force: 
    - in case of the high load 401 code are possible
    - not scalable - the biggest problem
* works from beginning

<BR>
<BR>
<BR>
<em>
We need something different... stickiness ?
</em>

</div>

<!--s-->

<!-- .slide: data-background="images/default-text-background-black-text-2.png" data-background-size="1920px" -->

<div id=left2-small>

#### Stickiness solution - cloudflared

* part of the apigw is side car, which connect directly to cf
* on cloudflare is possible to configure lb based on the user request
* pros:
    - supported solution!
* cons: 
    - sidecar overhead
    - expensive



<BR>
<BR>

<em>
unfortunately we didn't pass the perf test
</em>

</div>

<!--s-->


<!-- .slide: data-background="images/default-text-background-black-text-2.png" data-background-size="1920px" -->

<div id=left2-small>

#### Stickiness solution - haproxy in openshift

* haproxy is placed on the ocp
* stickiness based on the auth_token => user request from one user is placed on the same pods
* pros: cheap!
* cons: one hope plus overhead

<BR>
<BR>

<em>
unfortunately we didn't pass the perf test
</em>

</div>

<!--s-->

<!-- .slide: data-background="images/the-syncer-solution.png" data-background-size="1920px" -->

<div id=right3-small>

#### Final solution - syncer!

* developed in Go
* focused on the consistency instead of stickeness
* affect only login phase, all other requests work as before. No penalty for additional layer!
* customization of the js login process was necessary

</div>

<!--s-->
<!-- .slide: data-background="images/the-syncer-solution.png" data-background-size="1920px" -->

<div id=right3-small>

#### Syncer - How it works

* login process sends the origin and token key to the syncer in the end.
* syncer has info about running pods.
* it starts request in separate thread for each pods (except origin)
* (O) depends on the level of replication

</div>

<!--s-->
<!-- .slide: data-background="images/the-syncer-solution.png" data-background-size="1920px" -->

<div id=right3-small>

#### Syncer - How it works

* if the pod returns 401, repeat after some time
* if the pod returns 200, pod has a valid token
* when all threads finish with 200 => it's safe to return token to FE

</div>

<!--s-->

<!-- .slide: data-background="images/default-text-background-white-text-2.png" data-background-size="1920px" -->

<div id=left2>

#### Configuration

```yaml namespace: "apigwp-cz"
deployment: "ng-plus-apigw"
host: "api-apigwp-cz.t.dc1.cz.ipa.ifortuna.cz"
domain: "ifortuna.cz"
path: "/check"
port: "8080"
sync_timeout: 200
connection_timeout: 200
retries: 5
request_deadline: 1000
```

</div>

<!--s-->

<!-- .slide: data-background="images/default-text-background-white-text-2.png" data-background-size="1920px" -->

<div id=left2>

#### Apigw settings

```conf 
        location = /_synced {
            set $locid 0;
            auth_jwt "" token=$session_jwt;
            auth_jwt_key_request /_jwks_uri;
            limit_except GET {
                deny all;
            }         
            status_zone "FEG syncer";
            error_page 500 502 504 @oidc_error;
            js_content oidc.synced;
        }             
                      
        location = /_syncer {
            internal; 
            set $locid 0;
            auth_jwt_key_request /_jwks_uri;
            proxy_method GET;
            proxy_set_header X-Nginx-Origin $hostname;
            proxy_set_header X-Auth-Token $token_key;
            proxy_pass $ng_syncer_endpoint;
        }  
```

</div>

<!--s-->