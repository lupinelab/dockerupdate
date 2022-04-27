import argparse
from os import listdir
import subprocess
import docker as dkr

parser = argparse.ArgumentParser(description='Update docker images and rebuild containers')
parser.add_argument("-s", "--single", type=str, nargs=1, help="update single image/container")
args = parser.parse_args()

docker_client = dkr.from_env()


def update(docker):
    print(f"Stopping {docker} container:")
    stop = subprocess.run(["docker", "stop", docker], capture_output=True, text=True)
    print(stop.stdout)
    print(f"Removing {docker} container:")
    remove_container = subprocess.run(["docker", "rm", docker], capture_output=True, text=True)
    print(remove_container.stdout)
    print(f"Removing current {docker} image:")
    with open(f"/home/jedrw/dockercreate/{docker}", "r") as dockercreatefile:
        dockerregistry = (dockercreatefile.readlines()[-1].strip("\n").strip())
    remove_image = subprocess.run(["docker", "rmi", dockerregistry], capture_output=True, text=True)
    print(remove_image.stdout)
    print(f"Pulling latest {docker} image:")
    pull = subprocess.run(["docker", "pull", dockerregistry], capture_output=True, text=True)
    print(pull.stdout)
    print(f"Creating {docker} container:")
    create = subprocess.run(["sh", f"/home/jedrw/dockercreate/{docker}"], capture_output=True, text=True)
    print(create.stdout)
    print(f"Starting {docker} container:")
    start = subprocess.run(["docker", "start", docker], capture_output=True, text=True)
    print(start.stdout)


if args.single:
    update(args.single[0])
    print(f"{args.single[0]} status:")
    container = docker_client.containers.get(args.single[0])
    state = container.attrs["State"]
    print(state["Status"])
else:
    dockers = listdir("/home/jedrw/dockercreate")
    for docker in dockers:
        update(docker)
        print(f"{docker} status:")
        container = docker_client.containers.get(docker)
        state = container.attrs["State"]
        print(state["Status"])