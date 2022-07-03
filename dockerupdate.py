import argparse
from os import listdir, getlogin
import subprocess
from textwrap import dedent
import docker as dkr

parser = argparse.ArgumentParser(description='Update docker images or rebuild container(s)')
parser.add_argument("-i", "--image", type=str, nargs='?', const="all", help="update image and recreate container")
parser.add_argument("-c", "--container", type=str, nargs='?', const="all", help="recreate container")
args = parser.parse_args()
username = getlogin()
docker_client = dkr.from_env()
containers = listdir(f"/home/{username}/dockercreate")
builddir = listdir(f"/home/{username}/dockerbuild")


def remove_container(container):
    print(f"Stopping {container} container:")
    stop = subprocess.run(["docker", "stop", container], capture_output=True, text=True)
    print("Success")
    print(f"\nRemoving {container} container:")
    remove_container = subprocess.run(["docker", "rm", container], capture_output=True, text=True)
    print("Success")


def create_container(container):   
    print(f"\nCreating {container} container:")
    create = subprocess.run(["sh", f"/home/{username}/dockercreate/{container}"], capture_output=True, text=True)
    print(create.stdout)   
    print(f"Starting {container} container:")
    start = subprocess.run(["docker", "start", container], capture_output=True, text=True)
    print(get_status(container))


def update_image(container):
    print(f"\nRemoving current {container} image:")
    with open(f"/home/{username}/dockercreate/{container}", "r") as dockercreatefile:
        registry = (dockercreatefile.readlines()[-1].strip("\n").strip())
        imageid = subprocess.run(["docker", "images", "-q", registry], capture_output=True, text=True).stdout.strip("\n")
    print(f"{registry} - {imageid}:")
    remove_image = subprocess.run(["docker", "rmi", "-f", imageid], capture_output=True, text=True)
    print(remove_image.stdout)
    if container in builddir:
        print(f"Building {container} image:")
        build_image = subprocess.run(["docker", "build", f"/home/{username}/dockerbuild/{container}/", "-t", f"{registry}:latest"], capture_output=True, text=True)
        print(build_image.stdout)
        print(f"Pushing {container} image:")
        push_image = subprocess.run(["docker", "push", f"{registry}:latest"], capture_output=True, text=True)
        print(push_image.stdout)
    print(f"Pulling latest {container} image:")
    pull = subprocess.run(["docker", "pull", registry], capture_output=True, text=True)
    print(pull.stdout.strip())


def get_status(container):
    container = docker_client.containers.get(container)
    state = container.attrs["State"]["Status"]
    return state.capitalize()

if args.container:
    if args.container == "all":
        for container in containers:
            print(container.upper())
            print("=" * len(container.upper())) 
            remove_container(container)
            create_container(container)
            print("\n")
        print("Status Summary")
        print("--------------")
        for container in containers:
            print(f"{container + ':' :<30}:{get_status(container)}")
    else:
        print(args.container.upper())
        print("=" * len(args.container.upper()))
        remove_container(args.container)
        create_container(args.container)

elif args.image:
    if args.image == "all":
        for container in containers:
            print(container.upper())
            print("=" * len(container.upper())) 
            remove_container(container)
            update_image(container)
            create_container(container)
            print("\n")
        print("Status Summary")
        print("==============")
        for container in containers:
            print(f"{container + ':' :<30}:{get_status(container)}")
    else:
        print(args.image.upper())
        print("=" * len(args.image.upper()))
        remove_container(args.image)
        update_image(args.image)
        create_container(args.image)