import os

GO_PACKAGE_NAME = "polypheny.com/protos" # the package name for the proto files

PROTOS_PATH = "protos" # where proto files are

OUTPUT_PATH = "." # path for generated go code

with open("generate.sh", "w") as f:
    f.write("#!/bin/bash\n\n")

    script = f"protoc --proto_path={PROTOS_PATH} --go_out={OUTPUT_PATH} --go-grpc_out={OUTPUT_PATH} "

    names = []

    for protofile in os.listdir(PROTOS_PATH):
        if ".proto" in protofile: # you can never be too careful :)
            script += f" --go_opt=M{protofile}={GO_PACKAGE_NAME} --go-grpc_opt=M{protofile}={GO_PACKAGE_NAME} "
            names.append("protos/" + protofile)
    script += " ".join(names)
    f.write(script)
