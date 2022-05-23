import argparse
from os import listdir, getlogin
import subprocess
import docker as dkr

parser = argparse.ArgumentParser(description='Update docker images or rebuild container(s)')
parser.add_argument("-i", "--image", type=str, nargs='?', const="all", help="update image and recreate container")
parser.add_argument("-c", "--container", type=str, nargs='?', const="all", help="recreate container")
args = parser.parse_args()
username = getlogin()
docker_client = dkr.from_env()
dockers = listdir(f"/home/{username}/dockercreate")


def remove_container(docker):
    print(f"Stopping {docker} container:")
    stop = subprocess.run(["docker", "stop", docker], capture_output=True, text=True)
    print(stop.stdout)
    print(f"Removing {docker} container:")
    remove_container = subprocess.run(["docker", "rm", docker], capture_output=True, text=True)
    print(remove_container.stdout)


def create_container(docker):   
    print(f"Creating {docker} container:")
    create = subprocess.run(["sh", f"/home/{username}/dockercreate/{docker}"], capture_output=True, text=True)
    print(create.stdout)   
    print(f"Starting {docker} container:")
    start = subprocess.run(["docker", "start", docker], capture_output=True, text=True)
    print(start.stdout)


def update_image(docker):
    print(f"Removing current {docker} image:")
    with open(f"/home/{username}/dockercreate/{docker}", "r") as dockercreatefile:
        dockerregistry = (dockercreatefile.readlines()[-1].strip("\n").strip())
    remove_image = subprocess.run(["docker", "rmi", dockerregistry], capture_output=True, text=True)
    print(remove_image.stdout)
    if docker in listdir(f"/home/{username}/dockerbuild"):
        print(f"Building {docker} image:")
        build_image = subprocess.run(["docker", "build", f"/home/{username}/dockerbuild/{docker}/", "-t", f"{dockerregistry}:latest"], capture_output=True, text=True)
        print(build_image.stdout)
        print(f"Pushing {docker} image:")
        push_image = subprocess.run(["docker", "push", f"{dockerregistry}:latest"], capture_output=True, text=True)
        print(push_image.stdout)
    print(f"Pulling latest {docker} image:")
    pull = subprocess.run(["docker", "pull", dockerregistry], capture_output=True, text=True)
    print(pull.stdout)


def get_status(docker):
    print(f"{docker} status:")
    container = docker_client.containers.get(docker)
    state = container.attrs["State"]
    print(state["Status"])

if args.container:
    if args.container == "all":
        for docker in dockers:
            remove_container(docker)
            create_container(docker)
            get_status(docker)
    else:
        remove_container(args.container)
        create_container(args.container)
        get_status(args.container)

elif args.image:
    if args.image == "all":
        for docker in dockers:
            remove_container(docker)
            update_image(docker)
            create_container(docker)
            get_status(docker)
    else:
        remove_container(args.image)
        update_image(args.image)
        create_container(args.image)
        get_status(args.image)
