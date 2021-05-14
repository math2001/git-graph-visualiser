import json
import subprocess

# sparse: don't remove duplicate branches
lines = (
    subprocess.check_output(["git", "rev-list", "--all", "--children", "--sparse"])
    .decode("utf-8")
    .splitlines()
)


commits = {}

all_hashes = set()
all_children = set()  # all the hashes that are a child of another commit

for line in lines:
    (hash, *children) = line.split(" ")
    commits[hash] = {
        "message": subprocess.check_output(["git", "show", hash, "-s", "--format=%B"])
        .decode("utf-8")
        .splitlines()[0]
        .strip(),
        "children": list(
            sorted(children)
        ),  # sort (remember, these are commit hash) to get consistent graphs
    }
    all_hashes.add(hash)
    for child in children:
        all_children.add(child)

branches = {}
for line in subprocess.check_output(["git", "show-ref"]).decode("utf-8").splitlines():
    hash, ref = line.split(" ", 2)
    if ref.startswith("refs/heads/"):
        branches[ref.replace("refs/heads/", "")] = hash


print(
    json.dumps(
        {
            "commits": commits,
            "roots": list(all_hashes - all_children),
            "branches": branches,
            "HEAD": subprocess.check_output(["git", "symbolic-ref", "HEAD"])
            .decode("utf-8")
            .replace("refs/heads/", "")
            .strip(),
        },
    )
)