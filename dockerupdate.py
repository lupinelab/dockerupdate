from argparse import ArgumentParser
from os import listdir, getlogin, getcwd, chdir
from yaml import safe_load, YAMLError
from subprocess import run as subrun
from docker import docker_client as dkr

parser = ArgumentParser(description='Update docker images or rebuild container(s)')
parser.add_argument("-i", "--image", type=str, nargs='?', const="all", help="update image and recreate container")
parser.add_argument("-c", "--container", type=str, nargs='?', const="all", help="recreate container")
parser.add_argument("-b", "--build", type=str, nargs='?', const="", help="build image")
args = parser.parse_args()
username = getlogin()
docker_client = dkr.from_env()
containers = listdir(f"/home/{username}/docker")


def build_image(container):
    print(f"Building {container} image:")
    wd = getcwd()
    chdir(f"/home/{username}/docker/{container}/")
    build_image = subrun(["docker-compose", "build"], capture_output=True, text=True)
    print(build_image.stdout)
    chdir(wd)


def update_image(container):   
    print(f"Pulling {container} image:")
    wd = getcwd()
    chdir(f"/home/{username}/docker/{container}/")
    pull = subrun(["docker-compose", "pull"], capture_output=True, text=True) 
    print(f"{get_registry(container)}")
    print(f"Updating {container} container:")
    up = subrun(["docker-compose", "up", "--force-recreate", "-d"], capture_output=True, text=True)
    chdir(wd)
    print(f"{get_status(container)}\n")


def update_container(container):     
    print(f"Updating {container} container:")
    wd = getcwd()
    chdir(f"/home/{username}/docker/{container}/")
    up = subrun(["docker-compose", "up", "--force-recreate", "-d"], capture_output=True, text=True)
    chdir(wd)
    print(f"{get_status(container)}\n")


def get_registry(container):
    with open("docker-compose.yml", "r") as stream:
        try:
            compose = yaml.safe_load(stream)
            return compose['services'][container]['image']
        except yaml.YAMLError as exc:
            print(exc)

def get_status(container):
    container = docker_client.containers.get(container)
    state = container.attrs["State"]["Status"]
    return state.capitalize()

if args.container:
    if args.container == "all":
        for container in containers:
            print(container.upper())
            print("=" * len(container.upper())) 
            update_container(container)
        print("Status Summary")
        print("==============")
        for container in containers:
            print(f"{container}:{' ' * (30-len(container))}{get_status(container)}")
    else:
        print(args.container.upper())
        print("=" * len(args.container.upper()))
        update_container(args.container)

elif args.image:
    if args.image == "all":
        for container in containers:
            print(container.upper())
            print("=" * len(container.upper())) 
            update_image(container)
        print("Status Summary")
        print("==============")
        for container in containers:
            print(f"{container}:{' ' * (30-len(container))}{get_status(container)}")
    else:
        print(args.image.upper())
        print("=" * len(args.image.upper()))
        update_image(args.image)

elif not args.build:
    print("No image selected")

elif args.build:    
    print(args.build.upper())
    print("=" * len(args.build.upper()))
    build_image(args.build)