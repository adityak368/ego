#### Build and release Steps

1. Make changes and commit them

    ```
    git commit -m "Some Message"
    ```

2. Add one or more tags for that commit

    ```
    git tag -a broker/v1.0.9 -m "New Broker Release"
    git tag -a client/v1.0.5 -m "New Client Release"
    ```

    (2 tags are added for the commit)

3.  Push the commit with tags

    ```
    git push --tags origin
    ```

Adding tags to a specific commit

    git tag -a broker/v1.0.9 7cceb02 -m "Your message here"