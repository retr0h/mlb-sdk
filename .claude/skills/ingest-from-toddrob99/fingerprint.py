#!/usr/bin/env python3
"""Compute a canonical SHA-256 fingerprint for each endpoint defined in
toddrob99/MLB-StatsAPI's statsapi/endpoints.py.

This tool is consumed by the ingest-from-toddrob99 skill. Two reasons it
parses the file via Python's ast module rather than importing it:

1. statsapi/__init__.py imports `requests`, which we may not have available.
2. The fingerprint must be reproducible regardless of Python version, so we
   serialize each endpoint dict to canonical JSON (sorted keys, no spaces)
   before hashing — independent of dict iteration order.

Usage:
    python3 fingerprint.py <path/to/toddrob99/clone>
        → prints {endpoint_name: sha256_hex} as JSON to stdout.

The skill reads this output, compares it against manifest.json's
`upstreamSha` field per entry, and reports any drift since the last ingest.
"""

import ast
import hashlib
import json
import sys


def to_python(node, base_url):
    """Recursively turn an ast node into a plain Python value,
    resolving BASE_URL references."""
    if isinstance(node, ast.Dict):
        return {
            to_python(k, base_url): to_python(v, base_url)
            for k, v in zip(node.keys, node.values)
        }
    if isinstance(node, ast.List):
        return [to_python(e, base_url) for e in node.elts]
    if isinstance(node, ast.Tuple):
        return [to_python(e, base_url) for e in node.elts]
    if isinstance(node, ast.Constant):
        return node.value
    if isinstance(node, ast.Name):
        if node.id == "BASE_URL":
            return base_url
        if node.id in ("True", "False"):  # paranoia for older AST shapes
            return node.id == "True"
        raise ValueError(f"unexpected name reference: {node.id}")
    if isinstance(node, ast.BinOp) and isinstance(node.op, ast.Add):
        return to_python(node.left, base_url) + to_python(node.right, base_url)
    raise ValueError(f"unexpected node: {ast.dump(node)}")


def parse_endpoints(path):
    src = open(path).read()
    tree = ast.parse(src)

    base_url = None
    endpoints_node = None
    for node in tree.body:
        if not isinstance(node, ast.Assign):
            continue
        for target in node.targets:
            if not isinstance(target, ast.Name):
                continue
            if target.id == "BASE_URL" and isinstance(node.value, ast.Constant):
                base_url = node.value.value
            elif target.id == "ENDPOINTS":
                endpoints_node = node.value

    if base_url is None:
        raise SystemExit("BASE_URL not found in endpoints.py")
    if endpoints_node is None:
        raise SystemExit("ENDPOINTS dict not found in endpoints.py")

    return to_python(endpoints_node, base_url)


def fingerprint(entry):
    canonical = json.dumps(entry, sort_keys=True, separators=(",", ":"))
    return hashlib.sha256(canonical.encode("utf-8")).hexdigest()


def main():
    if len(sys.argv) != 2:
        sys.exit(f"usage: {sys.argv[0]} <path-to-toddrob99-clone-or-endpoints.py>")
    path = sys.argv[1]
    if not path.endswith("endpoints.py"):
        path = f"{path.rstrip('/')}/statsapi/endpoints.py"

    endpoints = parse_endpoints(path)
    out = {name: fingerprint(entry) for name, entry in endpoints.items()}
    print(json.dumps(out, indent=2, sort_keys=True))


if __name__ == "__main__":
    main()
