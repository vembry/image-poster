# challenge

## stories
- As a user, I should be able to: 
    - create posts with images (1 post - 1 image) 
    - set a text caption when I create a post
- As a user, I should be able to:
    - comment on a post
    - delete a comment (created by me) from a post
- As a user, I should be able to get the list of all posts along with the last 2 comments on each post

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