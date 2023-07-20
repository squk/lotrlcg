load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "test_lib",
    srcs = ["main.go"],
    importpath = "github.com/squk/lotrlcg/cmd/test",
    visibility = ["//visibility:private"],
    deps = [
        "@com_github_davecgh_go_spew//spew",
        "@com_github_mattn_go_sqlite3//:go-sqlite3",
    ],
)

go_binary(
    name = "test",
    embed = [":test_lib"],
    visibility = ["//visibility:public"],
)
