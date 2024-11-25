# challenge

## stories
- As a user, I should be able to: 
    - create posts with images (1 post - 1 image) - done
    - set a text caption when I create a post - done
- As a user, I should be able to:
    - comment on a post - done
    - delete a comment (created by me) from a post - done
- As a user, I should be able to get the list of all posts along with the last 2 comments on each post - done

## functional requirements
- RESTful Web API (JSON) - done
- Maximum image size - 100MB - done
- Allowed image formats: .png, .jpg, .bmp. - done
- Save uploaded images in the original format - done
- Convert uploaded images to .jpg format and resize to 600x600 - done
- Serve images only in .jpg format - done, actually i provide both url to original and converted
- Posts should be sorted by the number of comments (desc) - done
- Retrieve posts via a cursor-based pagination - done

## non-functional requirements

- Maximum response time for any API call except uploading image files - 50 ms
- Minimum throughput handled by the system - 100 RPS
- Users have a slow and unstable internet connection


# what i did

## tech spec

- `golang` v1.23.3 as webapp
- `postgres` @ latest version for storage
- `docker` v27.0.3-1
- `localstack` v4.0.2 as aws alternative for local development since i uses S3 and SQS.
    - s3 to store the images uploaded
    - sqs as worker queue to transform file
    - ref: [link](https://www.localstack.cloud/)
- ubuntu 24.04.1 LTS
- vscode optional
- vscode's Dev Container, optional
    - ref to vscode's Dev Container documentation: [link](https://code.visualstudio.com/docs/devcontainers/containers)
    - this help isolate the development in container, rather than requiring dev to setup everything independently

## entity
I may have shoot myself in the foot here, I'm attempting to create a cyclic structure. So essentially both post and comment behave like post. The key difference is how we present them. The structure are stored on `post_structures` which contain only `post_id` and `parent_post_id`. `post_structures` entries with `parent_post_id` = NULL will be counted as Posts, and for entries with `parent_post_id` != NULL will be counted as Comments. I take inspiration from reddit.

## flow

### `create post`
1. user send request
    - right now i mock user authentication by defining x-user-id directly
2. app receive request
    1. validate file size < 100MB
    2. validate data type, only allows jpg, jpeg, png, bmp
3. upload file to s3
4. create post entry to `posts` table
5. create post-structure entry to `post_structures` table
6. enqueue task to sqs for image transform as per noted on challenge

### `transform image`
1. consume task enqueued by `create post` flow
2. retrieve post detail from `posts` table
3. retrieve image from s3
4. transform image
5. upload it to s3 as new file
6. update image details on post data
7. post data on `posts` table


### `get list of posts`
1. user send request
2. app receive request
3. retrieve posts 
    1. based on limit/page combo
    2. ordered by comment count
4. retrieve comments of those posts
    1. limit to only 2 latest comments per posts

### `post comment`
1. user send request
2. app receive request
3. validate post(comment target) existence
    1. if not exists, then do early return
4. create comment entry to `posts` table
5. create comment entry to `post_structures` table with `parent_post_id` pointing to the post it comments

### delete comment
1. user send request
2. app receive request
3. validate comment ownership
    1. if request attempts to delete comment belong to someone else, then do early return
4. soft delete comments on `posts`


## project structure


I'm taking Lego blocks as inspiration on building the web-app. Essentially every block can be stacked together in order for us to create something, in this case a web-app. These "building blocks" are defined explicitly and scaffolded together in the `main.go`, which will also act as the web-app's orchestrator. This orchestration covers:
1. starting the web-app, 
2. waiting for exit signal, and 
3. shuting-down the web-app.  

### `./.aws`
contain mock aws config and creds for localstack to work 

### `./.devcontainer`
contains vscode's dev-container configuration

### `./.docker`
contains scripts to help setup the local cluster

### `./app-go`
contains the web-app

#### `./app-go/configs`
contains env vars 

#### `./app-go/internal`
contains implementation for the app-go

#### `./app-go/internal/app`
contain basic helper for general purpose usage like config, logger, etc

#### `./app-go/internal/clients`
contain outbound client implementation, right now it's not really consistent as i still have definition all over the place. But ideally outbound client should be located here.

#### `./app-go/internal/models`
contain common data transfer objects / models

#### `./app-go/internal/modules`
contain domains which has it's own service and repo, 

#### `./app-go/internal/servers`
contain implementatin to serve the app

#### `./app-go/internal/workers`
contain worker implementation, right now it only contain sqs workers

#### `./app-go/main.go`
the app's orchestrator

### `./compose.yml`
compose file to run local cluster in docker

### `./makefile`
to assists the setup

### `./postman`
containing curls for the app-go http

