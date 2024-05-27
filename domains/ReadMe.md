1. To create a library in golang, use go_library module
    - Create a BUILD file and insert following contents to it
        ```
        go_library(
            name = "go_simple",                                                                        Target Name
            srcs = ["simple.go"],                                                                      Source File 
            importpath = "github.com/Akhilesh53/ArtHub/Projects/go_folder/go_simple",                  Import Path (most important, without this we cannot use the library) 
            visibility = ["//visibility:public"],                                                      Visibility of the library (check out for more labels) 
        )
        ```

    - To test a functionality, Use go_test module. Add the same contents in the same build file.
        ```
        go_test(
            name = "simple_test",                    Target Name
            srcs = ["simple_test.go",],              Source File (Where tests are defined) 
            embed = [":go_simple"],                  embed is like, which source dependency to use
        )
        ```

2. To use the created library in another project or folder, Add the folder path in dependencies (deps). and import the library using the importpath specified (i.e  github.com/Akhilesh53/ArtHub/Projects/go_folder/go_simple")
    - 
    ```
    go_binary(
        name="go_web",
        srcs = ["main.go"],
        deps = [
            "//Projects/go_folder/go_simple",          Path to library folder
        ]
    )
    ```
3. There are many ways to import external libraries but the straight forward approach is to include the repository in the workspace.

    In the WORKSPACE.bazel file, at line 31 we have imported an external golang library using go_repository module
      ```
      go_repository(
        name = "com_github_gorilla_mux",                            Name can be anything, By convention it should be the reverse mapping of url
        importpath = "github.com/gorilla/mux",                      Path from where library should be imported  
        tag = "v1.8.0",                                             Latest version tag
        #sum
        #version                                                    For production, better is to use sum and version
    )
    ```

    Now, How do we get to know what should we write in dependencies to import the library.

    Run bazel 
    ```
    query @com_github_gorilla_mux //...
    ```

    From the command, we get an output where this library is used as a dependency and it wil give us an import path.

    ```
    go_binary(
        name="go_web",
        srcs = ["main.go"],
        deps = [
            "//Projects/go_folder/go_simple",
            "@com_github_gorilla_mux//:mux"                            Import Path generated from Output
        ]
    )
    ```

    Build and Run the files.
    - To Build - 
    ```
    build bazel run //Projects/go_folder/go_web:go_web
    ```

    - To Run  - 
    ```
    bazel run //Projects/go_folder/go_web:go_web
    ```
