# Monorepo
Monorepo SetUp Using Bazel

#
### To build the BUILD files of all the projects
   Run `bazel build //...`

#
### To build BUILD file of any particular project
   Run `bazel build //path/to/project: target-name`

#
### To Run any target file
   Run `bazel run //path/to/project: target-name`

#
### To Test all test files
   Run `bazel test //...`

#
### To test any target file
   Run `bazel test //path/to/test: target-name`

#
### To get the URL / string for dependencies to import external library
   Run `bazel query @com_github_gorilla_mux //...`
    `bazel query @target_name //...`

#
### To generate BUILD fles using gazelle
   Run `bazel run //:gazelle`

#
### To update the dependency/repos from mod files using gazelle
   Run `bazel run //:gazelle-update-repos`

`NOTE:` Prefer not to use gazelle. It will regenerate BUILD files which may have wrong configurations

#

- `//` - means the current directory (where workspace is located)
- `...` - means all files

#

### TO SETUP DOCKER
- Include the docker rules in the `WORKSPACE` file.


  ```
  http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "b1e80761a8a8243d03ebca8845e9cc1ba6c82ce7c5179ce2b295cd36f7e394bf",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.25.0/rules_docker-v0.25.0.tar.gz"],
   )

   load(
      "@io_bazel_rules_docker//repositories:repositories.bzl",
      container_repositories = "repositories",
   )
   container_repositories()

   load("@io_bazel_rules_docker//repositories:deps.bzl", container_deps = "deps")

   container_deps()
  ```

- To load any image, load the image using the beelow mentioned code in `WORKSPACE` file.
  ```
  load(
    "@io_bazel_rules_docker//python3:image.bzl",
    _py_image_repos = "repositories",
   )

   _py_image_repos()
  ``` 
  This code loads the python3 image from docker rules.


- To use this image in any application, write the following code in `BUILD` file.
  ```
  load("@io_bazel_rules_docker//python3:image.bzl", py_image = "py3_image")                       language specific image rule

  py_image(
    name = "samplepython_image",                                                                  name of the target image
    srcs = ["app.py"],                                                                            source code to be converted to image
    main = "app.py",                     
    deps = [                                                                                      any dependencies, if there is 
        "//projects/python_folder/sample_library:calculator",
        requirement("Flask"),
    ],
  )
  ```

  In the above mentioned code, we have used languge specific image. We can use `container_image` in case we dont want to use any language specific image.

  #
  ## NOTE:
  These images not directly run on MacOS. An error will pop up while trying to run this docker image.
  ```
  The requested image's platform (linux/amd64) does not match the detected host platform (linux/arm64/v8) 
  ``` 

  To resolve this error, we need to specify the platform for which image should run.
  Add the below command to `.bazlrc` file in root directory
  ```
  build --@io_bazel_rules_docker//transitions:enable=no --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64
  ```

  #

- To build the file 
  
  RUN ```bazel build //...```

- Or to build only specified image

  RUN ```bazel build //projects/python_folder/python_app:samplepython_image```

These build command will assign an id to the image. 

- To run the image 

  RUN ```docker  run -p 8087:8087 bazel/projects/python_folder/python_app:samplepython_image```

