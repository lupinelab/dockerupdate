import argparse
from os import listdir
import subprocess


parser = argparse.ArgumentParser(description='Update docker images and rebuild containers')
parser.add_argument("-s", "--single", type=str, nargs=1, help="update single image/container")
args = parser.parse_args()


def update(docker):
    print(f"Stopping {docker} container:")
    subprocess.run(["docker", "stop", docker])
    print(f"Removing {docker} container:")
    subprocess.run(["docker", "rm", docker])
    print(f"Removing current {docker} image:")
    with open(f"/home/jedrw/dockercreate/{docker}", "r") as dockercreatefile:
        dockerregistry = (dockercreatefile.readlines()[-1].strip("\n").strip())
    subprocess.run(["docker", "rmi", dockerregistry])
    print(f"Pulling latest {docker} image:")
    subprocess.run(["docker", "pull", dockerregistry])
    print(f"Creating {docker} container:")
    subprocess.run(["sh", f"/home/jedrw/dockercreate/{docker}"])
    print(f"Starting {docker} container:")
    subprocess.run(["docker", "start", docker])
    print(f"{docker} status:")
    subprocess.run("docker ps --filter name=" + docker + " --filter status=running", shell=True)


if args.single:
    update(args.single[0])
else:
    dockers = listdir("/home/jedrw/dockercreate")
    for docker in dockers:
        update(docker)