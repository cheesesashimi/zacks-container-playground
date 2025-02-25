# containerfile

While there is prior art for parsing an abstract syntax tree (AST) of Containerfiles, the reverse is not (yet) possible. This package aims to provide helpers for programmatically generating a Containerfile using some higher-level abstractions and primitives. Rather than use a Go template or other difficult-to-reason about ways of constructing a Containerfile, one can instantiate the structs contained within this package. Right now, no validation is performed however that can be added in the future.

For right now, these solely perform string interpolation and concatenation to construct individual Containerfile statements and directives. In the future, this could interact directly with container image build APIs to more directly perform the requested actions.
