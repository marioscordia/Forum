# Web Forum

This is a web forum application, which provides certain functionalities of social network like creating posts with image uploading, liking/dislikng, commenting, notifications and deleting. The project is built using Golang and SQLite.

It implements 4 types of users:

1.  Guests

    - These are unregistered-users that can neither post, comment, like or dislike a post. They only have the permission to see those posts, comments, likes or dislikes.

2.  Users

    - These are the users that will be able to create, comment, like or dislike posts and also delete their own posts and comments.

3.  Moderators

    - Moderators, are users that have a granted access to special functions :
      - They are responsible for which posts to be published
      - They are be able to monitor the content in the forum by reporting post to the admin

4.  Administrators

    - Users that manage the technical details required for running the forum. This user must be able to :
      - Promote or demote a normal user to, or from a moderator user.
      - Receive reports from moderators. If the admin receives a report from a moderator, he can respond to that report
      - Delete posts and comments

## Note

- There is only one administrator, which is created at initialization of the database.
  Login : **meduza@gmail.com**
  Password : **Cheburek**
- Images are not saved in database itself, since it may slow down the database operations. Instead, images are saved in **./internal/store/img/** folder.

## Objectives

The objectives of this project are:

- Allow users to communicate between each other through posts and comments.
- Associate categories to posts.
- Allow users to like and dislike posts and comments.
- Implement a filter mechanism to filter posts by categories, created posts, and liked posts.
- Use SQLite to store the data for the application.
- Implement user authentication and sessions using cookies.
- Follow good coding practices and handle all sorts of errors.

## Installation

To run this project, you need to have Docker installed on your machine. Then follow these steps:

1. Clone this repository.
2. Navigate to the project directory.
3. Run the following command to build the Docker image:

```bash
make build
```

4. Run the following command to start the Docker container:

```bash
make run
```

5. Open your web browser and go to http://localhost:{port number}.

If you want to run without Docker:

```bash
go run /cmd/web/*
```

## Usage

To use the web forum application, follow these steps:

1. Register a new user by providing your email, username, and password.
2. Log in to the application using your email and password.
3. Create a post by providing a title, content, and one or more categories.
