load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

package(
    default_visibility = ["//visibility:public"],
)

go_library(
    name = "beornextract_lib",
    srcs = ["main.go"],
    data = [
        "//third_party/data/hallofbeorn/data:Bot.Cards.json",
        "//third_party/data/hallofbeorn/data:Export.Cards.json",
        "//third_party/data/hallofbeorn/data:cycles.json",
        "//third_party/data/hallofbeorn/data:ringsdb.json",
    ],
    importpath = "github.com/squk/lotrlcg/src/cmd/beornextract",
    deps = [
        "@com_github_grokify_html_strip_tags_go//:html-strip-tags-go",
        "@com_github_jessevdk_go_flags//:go-flags",
        "@com_github_satori_go_uuid//:go_uuid",
        "@com_github_squk_lotrlcg//cmd/beornextract/types:go_default_library",
        "@io_bazel_rules_go//go/tools/bazel:go_default_library",
    ],
)

go_binary(
    name = "beornextract",
    embed = [":beornextract_lib"],
)
