load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

package(
    default_visibility = ["//visibility:public"],
)

go_library(
    name = "beornextract_lib",
    srcs = ["main.go"],
    data = [
        "//cmd/beornextract/data:Bot.Cards.json",
        "//cmd/beornextract/data:Export.Cards.json",
    ],
    importpath = "github.com/squk/lotr/cmd/beornextract",
    deps = [
        "//cmd/beornextract/types",
        "@com_github_grokify_html_strip_tags_go//:html-strip-tags-go",
        "@com_github_jessevdk_go_flags//:go_default_library",
        "@io_bazel_rules_go//go/tools/bazel:go_default_library",
    ],
)

go_binary(
    name = "beornextract",
    embed = [":beornextract_lib"],
)
