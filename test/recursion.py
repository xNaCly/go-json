import os

depths = {
    "1K": 1_000,
    "10K": 10_000,
    "100K": 100_000,
    "1M": 1_000_000,
    "10M": 10_000_000,
}

for depth_name, depth in depths.items():
    print(f"Generating {depth} depth object")

    json_string = "{"

    for _ in range(1, depth):
        json_string += '"next":{'

    json_string += '"next":null'

    for _ in range(depth):
        json_string += "}"

    file_name = f"{depth_name}_recursion.json"
    file_path = os.path.join(os.path.dirname(__file__), file_name)

    with open(file_path, 'w') as f:
        f.write(json_string)

    print(f"File for depth {depth} saved as {file_name}")

