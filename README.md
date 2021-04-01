# How deploy project locally to kubernetes

0. Ensure you have Docker Desktop running.
1. Use `minikube start` to deploy single node cluster locally on your computer.
2. We will need to change the configuration for minikube to support ingress (ingress-nginx controller) which enables public access to our containers:

    ```
    $ minikube config set vm-driver hyperkit
    $ minikube delete
    $ minikube start
    $ minikube addons enable ingress
    ```

3. Since our server image is hosted on a private registry on Docker Hub, we need to create a secret for authentication (refer to `deployments/backend-deployment.yaml` `imagePullSecrets`) We will name our secret `regcred`:

    ```
    kubectl create secret docker-registry regcred --docker-server=<your-registry-server> --docker-username=<your-name> --docker-password=<your-pword> --docker-email=<your-email>
    ```

4. Now we can begin deployment of our containers and services:

    ```
    # Run this in root folder which will run all .yaml files in deployments folder
    $ kubectl apply -f deployments
    ```

5. View the status of all pods, services and deployments:
    
    ```
    $ kubectl get all
    ```

6. Once successfully deployed, you can access it locally:
    
    ```
    $ minikube ip
    ```
    This will return an external IP address that you can paste in your browser.

### Cleaning Up

It is important to clean up the resources you have deployed:
```
# At root folder
$ kubectl delete -f deployments
``` 

# How to update docker images

When deploying our application to a kubernetes cluster, our `.yaml` files will pull the docker images from our repositories from Docker Hub. To update images:
1. Make sure you are logged into Docker Hub:

    ```
    $ docker login
    ```
2. Build the docker image for either the server or client:

    ```
    $ docker build -t sync-server:v0.0 -f ./build/Dockerfile .
    ```

    - The `-t` argument tags the image. By convention, the label after the colon is used to specify the image version.

    - The `-f` argument specifies the path of the Dockerfile which can be useful if the dockerfile is not stored in the root folder of the application but you want the context to be at the root.

    This build will be saved locally.

3. Tag the image to prepare for it to be pushed to the repository:

    ```
    $ docker image tag sync-server:v0.0 <YOUR-DOCKER-ID>/sync-server:v0.0
    ```
4. Push to the repository:
    
    ```
    $ docker image push <YOUR-DOCKER-ID>/sync-server:v0.0
    ```