wrk.method = "POST"
wrk.body   = 'query "\\u4e16" foo { name model }'
wrk.headers["Content-Type"] = "application/graphql"
