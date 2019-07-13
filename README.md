# wsf-svc-example
service example base web server framework

# Default configuration
-- bin
   -- wsf-svc-example.exe  // service application
-- cfg
   -- wsf-svc-example.json // config file
-- crt
   -- server.pfx           // cirtificate for https server
-- site
   -- root                 // folder of root site
   -- doc                  // folder of api document site
   -- opt                  // folder of service administration site
   -- webapp               // root folder of web application site
   
   
# Install
wsf-svc-example.exe -install // install as a service
wsf-svc-example.exe -start   // run the service
wsf-svc-example.exe -stop    // stop the service
wsf-svc-example.exe -help    // display usage information
