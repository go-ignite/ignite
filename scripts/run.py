import string
import random
import socket
import commands
import sys

def get_port(port):
    address = '127.0.0.1'
    s = socket.socket()
    print "Attempting to connect to {} on port {}".format(address, port)
    try:
        s.connect((address, port))
        return True
    except socket.error, e:
        return False

def check_docker(name):
    result = commands.getstatusoutput('docker ps -a|grep {}'.format(name))
    if 'Up' in result[1]:
	return True
    return False

if __name__ == '__main__':

    pwd_len = 15
    random_pass = ''.join(random.choice(string.ascii_lowercase + string.ascii_uppercase + string.digits) for _ in range(pwd_len))
    print "New random password is: {}".format(random_pass)

    for port in range(4000,5001):
        if not get_port(port):
            print "Port {} is available!".format(port)
	    name = sys.argv[1]

            result = commands.getstatusoutput('docker run -d --name={} -p {}:{} --restart=always registry.xiaozhou.net/ssmu:latest -p {} -k {}'.format(name, port, port, port, random_pass))
	    if result[0] == 0:
		print "Docker created..."

	    if check_docker(name):
	        print "Docker: {} is running!".format(name)
		print "Port is: {}".format(port)
	 	print "Password: {}".format(random_pass)
	    break
