# openaiproxy

Here is a "magic trio" making EventSource working through Nginx:

proxy_set_header Connection '';
proxy_http_version 1.1;
chunked_transfer_encoding off;
Place them into location section and it should work.

You may also need to add

proxy_buffering off;
proxy_cache off;
 
