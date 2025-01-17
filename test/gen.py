import os
from os.path import exists
import math

sizes =[1,5,10]

line = """\t{
        "key1": "value",
        "array": [],
        "obj": {},
        "atomArray": [11201,1e112,true,false,null,"str"]
    }"""

def write_data(size: int): 
    name = f"{size}MB.json"
    if not exists(name):
        with open(name, mode="w", encoding="utf8") as f:
            f.write("[\n")
            size = math.floor((size*1000000)/len(line))
            f.write(",\n".join([line for _ in range(0, size)]))
            f.write("\n]")

[write_data(size) for size in sizes]

depths = {
    "1K": 1_000,
    "10K": 10_000,
    "100K": 100_000,
    "1M": 1_000_000,
    "10M": 10_000_000,
}

for depth_name, depth in depths.items():
    print(f"Generating {depth} depth object")

    json_parts = ["{"]

    for _ in range(1, depth):
        json_parts.append('"next":{')

    json_parts.append('"next":null')

    for _ in range(depth):
        json_parts.append("}")

    json_string = "".join(json_parts)

    file_name = f"{depth_name}_recursion.json"
    file_path = os.path.join(os.path.dirname(__file__), file_name)

    with open(file_path, 'w') as f:
        f.write(json_string)

    print(f"File for depth {depth} saved as {file_name}")
