import urllib2
import os


def make_executable(path):
    mode = os.stat(path).st_mode
    mode |= (mode & 0o444) >> 2
    os.chmod(path, mode)


if __name__ == '__main__':
    response = urllib2.urlopen('https://github.com/thylong/ian/releases/download/ian-v0.2.4/ian-darwin-amd64')

    data = response.read()
    path = "/usr/local/bin/ian"

    directory = os.path.dirname("/usr/local/bin/")
    if not os.path.exists(directory):
        os.makedirs(directory)

    file_ = open(path, 'w+')
    file_.write(data)
    file_.close()

    make_executable(path)
