# Description
This is a simple single-paged web forum application built with Go, javascript, SQLite, and Docker. It allows users to communicate, create posts, comment, like/dislike content, and filter posts by categories or user activity.

# Features
- **user Authentication**
    - uesrs can register and create new accounts.
    - session management using sessions and cookies.
    - registered users can like, dislike, comment and create posts.
    - non-registered users can only view the content of the forum.
    - types of users:
        - Guest users: are non-logged in users, they have limited access to the forum which is restricted to only viewing the content of the forum.
        - Normal users: are logged in users, they can post, comment, like and dislike.
        - Moderator users: are logged in users with all the privileges of Normal users, with addition to their ability of deleting posts in the forum or report them to the administirators.
        - Administrator users: are logged in users with unlimited privileges, they can:
            - promote Normal users to moderators or demote moderator users to normal users.
            - Receive reports from moderators. If the admin receives a report from a moderator, he can respond to that report.
            - Delete posts and comments.
            -  manage categories by addind and deleting them.
- **posts and comments**
    - posts can be associated with categories
    - images can be added to a post
    - posts can be commented by users
- **likes and dislikes**
    - users can like posts & comments
    - when a non-registered user tries to like, they'll be redirected to the login page
- **filters**
    - in catigories page the user can filter posts by their category
- **security**
    - reate limiting is applied to the forum:
        - if the user is logged in: the user is going to be blocked by the server.
        - if the user is not logged: the IP address is going to be blocked by the server.
    - HTTPS: the forum uses HTTPS(Hyper Text Transfer Protocol Secure) for a secure connection.
    - password hashing: using the bcrypt lib to store the password hashes for better user security.

# How to use
1. **Clone the repository:**
   ```bash
   git clone https://learn.reboot01.com/git/musabt/forum.git
   cd forum
   ```
2. **Build the Docker image:**
    ```bash
    docker build -t forum-app .
    ```
3. **Run the application in a Docker container:**
    ```bash
    docker run -p 443:443 forum-app
    ```
4. **Access the forum:**
    Open your browser and navigate to https://localhost

# Authors
[@musabt AKA:MAISTRY](https://learn.reboot01.com/git/musabt)

[@mmahmooda AKA:KASIKO](https://learn.reboot01.com/git/mmahmooda)

[@alimarhoon AKA:NBX](https://learn.reboot01.com/git/alimarhoon)

[@malmannai AKA:mvniraa](https://learn.reboot01.com/git/malmannai)
