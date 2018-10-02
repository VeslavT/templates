from collections import defaultdict

graph = defaultdict(set)


def process_dfs():
    result = defaultdict(set)
    for key in sorted(list(graph.keys())):
        visited = dfs(key, set())
        visited.remove(key)
        result[key] = visited
        for ver in visited:
            if ver in graph:
                del graph[ver]


def dfs(vertex, visited):
    visited.add(vertex)
    if vertex in graph:
        for w in graph[vertex]:
            if w not in visited:
                dfs(w, visited)
    return visited