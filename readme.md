# challenge

## stories
- As a user, I should be able to create posts with images (1 post - 1 image)
- As a user, I should be able to set a text caption when I create a post
- As a user, I should be able to comment on a post
- As a user, I should be able to delete a comment (created by me) from a post
- As a user, I should be able to get the list of all posts along with the last 2 comments on each post

## functional requirements
- RESTful Web API (JSON)
- Maximum image size - 100MB
- Allowed image formats: .png, .jpg, .bmp.
- Save uploaded images in the original format
- Convert uploaded images to .jpg format and resize to 600x600
- Serve images only in .jpg format
- Posts should be sorted by the number of comments (desc)
- Retrieve posts via a cursor-based pagination

## non-functional requirements
- Maximum response time for any API call except uploading image files - 50 ms
- Minimum throughput handled by the system - 100 RPS
- Users have a slow and unstable internet connection