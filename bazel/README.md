# Bazel
Followed this tutorial:

Didn't complete this...

```sh
brew install bazel

cd java-tutorial
bazel build //:ProjectRunner


bazel-bin/ProjectRunner
```

https://docs.bazel.build/versions/3.6.0/tutorial/java.html



The `WORKSPACE` file, which identifies the directory and its contents as a Bazel workspace and lives at the root of the project’s directory structure,
One or more `BUILD` files, which tell Bazel how to build different parts of the project. (A directory within the workspace that contains a BUILD file is a package. You will learn about packages later in this tutorial.)
