wrk.method = "POST"
wrk.body   = "query foo { name model ! $ ( ) . : = @ [ ] { | } }"
wrk.headers["Content-Type"] = "application/graphql"
