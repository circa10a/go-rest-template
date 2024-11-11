docker_build('circa10a/go-rest-template', '.', dockerfile='./Dockerfile')
k8s_yaml(listdir('./deploy/k8s'))
k8s_resource('go-rest-template', port_forwards=8080)