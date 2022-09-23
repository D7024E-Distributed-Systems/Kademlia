Run in terminal:
docker exec -it <container name> bash
example: docker exec -it kademlia_kademliaNodes_1 bash

in the new terminal run:
reptyr 1


If you run this command and then start the docker containers you then can access those from the cmd go run src/main.go instead of having to docker exec in.
systemctl --user stop docker-desktop