package(
    default_visibility = ["//visibility:public"],
)

load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/squk/lotrlcg
# gazelle:build_file_name BUILD
gazelle(name = "gazelle")

gazelle(
    name = "gazelle-update-repos",
    args = [
        "-from_file=go.mod",
        "-to_macro=deps.bzl%go_dependencies",
        "-prune",
        "-build_file_proto_mode=disable_global",
    ],
    command = "update-repos",
)

alias(
    name = "beornextract",
    actual = "//cmd/beornextract:beornextract",
)
