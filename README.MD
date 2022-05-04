# forum

## To run the project locally

Run the following commands:
For building an exec file: **make build**
For running the project: **make run** or **go run main.go**
For creating an image and running container: **make docker**

## Admin Role account

    Login: admin
    Password: admin

## Description
**MVC** pattern was used to structure the forum

This web forum allows :

* communication between users
* associating categories to posts
* liking and disliking posts and comments
* filtering posts
* uploading images to posts
* updating/deleting posts and comments
* shows the user's activity page
* notifications
* separation users in roles: guests, users, moderators, admins


### Used technologies and packages
* SQLite
* Docker
* Standard go packages are used.
* github.com/mattn/go-sqlite3
* golang.org/x/crypto/bcrypt
* github.com/satori/go.uuid


### Authentication

Clients can **register** as a new user on the forum, by inputting their credentials.

Cookies was used as a way to authenicate users. 

User registration:
* Asks for email
    * When the email is already taken returns an error response
* Asks for username
* Asks for password
    * The password is encrypted when stored

### Communication

In order for users to communicate between each other, they can create posts and comments.

* Only registered users are able to create posts and comments.
* When registered users are creating a post they can associate one or more categories to it.
* The posts and comments are be visible to all users (registered or not).
* Non-registered users are only able to see posts and comments.

### Likes and Dislikes

Only registered users are able to like or dislike posts and comments.

The number of likes and dislikes are visible by all users (registered or not).

### Filter

Filter mechanism allows users to filter the displayed posts by :

* categories
* created posts
* liked posts

The last two are only available for registered users and refers to the logged in user.

### Image-Upload

Registered users have the possibility to create a post containing an image as well as text.

* When viewing the post, users and guests could see the image associated to it.

There are several extensions for images like: JPEG, PNG, GIF.

The max size of the images to load can be 20 mb. If there is an attempt to load an image greater than 20mb, an error message informs the user that the image is too big.

### Advanced-features

* Notifications are triggered when users posts are :
    * liked/disliked
    * commented

* Activity page tracks the user own activity:
    * Shows the user created posts
    * Shows where the user left a like or a dislike
    * Shows where and what the user has been commenting

* Users are able to to Edit/Remove posts and comments

### Moderation
There are four types of users :

#### Guests

* These are unregistered-users that can neither post, comment, like or dislike a post. They only have the permission to see those posts, comments, likes or dislikes.

#### Users

* These are the users that are able to create, comment, like or dislike posts.

#### Moderators

* Moderators are users that have a granted access to special functions :
    * They are able to monitor the content in the forum by deleting or reporting post to the admin
* To create a moderator the user should request an admin for that role

#### Administrators

* Users that manage the technical details required for running the forum. This user is able to :
    * Promote or demote a normal user to, or from a moderator user
    * Receive reports from moderators. If the admin receives a report from a moderator, he can respond to that report
    * Delete posts and comment
    * Manage the categories, by being able to create and delete them



