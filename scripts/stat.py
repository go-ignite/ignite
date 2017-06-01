import subprocess
import sys
import re 

def get_bandwidth(name):
    try:
        cmd = 'docker exec -it {} /sbin/ifconfig eth0'.format(name)
        result = subprocess.check_output(cmd, shell=True).decode('utf-8')

        print(result)
        match = re.search(r'TX bytes:\d+ ', result)
        raw = float(match.group(0).split(':')[1])
        bandwidth = raw/1024/1024

        if bandwidth > 1024:
            print('Bandwidth is: {} / {:.2f} GB'.format(raw, float(bandwidth)/1024))
        else:
            print('Bandwidth is: {} / {:.2f} MB'.format(raw, float(bandwidth)))
    except subprocess.CalledProcessError as ex:
        print('Failed to get docker bandwith, {} {}'.format(ex.returncode, ex.output))

if __name__ == '__main__':
    print('Running...')
    get_bandwidth(sys.argv[1])
