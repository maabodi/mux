# Backend-API
  
this app need bucket storage for images or file.  
i'm using minio bucket so, download minio server.
> ### LINUX  
``` https://min.io/download#/linux ```
  
after download, just running minio server.
> sudo \ \
MINIO_ROOT_USER=admin \ \
MINIO_ROOT_PASSWORD=admin123 \ \
./minio server {directory_for_storage} \ \
--console-address ":9001"
  
to access minio.
> ``` http://localhost:9001 ``` 
  
after minio running, we can start the app.  
app will running on port : ``` 1323 ```   
to run app : ``` go run main.go ```  
to access swagger : ``` http://localhost:1323/swagger/index.html ```