# Docker Practical examples and Commands cheatsheet
## Docker commands
 ```bash
 docker build -t cs_portfolio:0.1 .        # build the docker image for the project
 docker run cs_portfolio:0.1               # run the container without port won't make the website available
 docker run -p 80:80 cs_portfolio:0.1      # run the app with the same port of the container
 docker run -p 7785:80 cs_portfolio:0.1    # we can change the port on our local machine not the container side
 docker login                              # login to docker hub
 docker tag cs_portfolio:0.1 maziar/cs_portfolio:0.1            # tag the image
 docker push maziar/cs_portfolio:0.1                            # push the docker image to the docker hub remote repository
 ```
