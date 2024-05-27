load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_binary(
    name = "sample_app",
    embed = [":sample_app_lib"],
    visibility = ["//visibility:public"],
)

go_library(
    name = "sample_app_lib",
    srcs = ["main.go"],
    importpath = "github.com/monorepo/domains/sample_app",
    visibility = ["//visibility:private"],
    deps = [
        "//domains/sample_library",
        "@com_github_gorilla_mux//:mux",
    ],
)
