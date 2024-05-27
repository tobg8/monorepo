load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "sample_library",
    srcs = ["sample_library.go"],
    importpath = "github.com/monorepo/domains/sample_library",
    visibility = ["//visibility:public"],
)

go_test(
    name = "sample_test",
    srcs = ["sample_test.go"],
    embed = [":sample_library"],
)

go_test(
    name = "sample_library_test",
    srcs = ["sample_test.go"],
    embed = [":sample_library"],
)
