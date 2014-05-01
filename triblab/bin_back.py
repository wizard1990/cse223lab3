import subprocess
import time

global server_name
global running_server
server_name = []
server_name.append("localhost:10000")
server_name.append("localhost:10001")
server_name.append("localhost:10002")
server_name.append("localhost:10003")
server_name.append("localhost:10004")
server_name.append("localhost:10005")
server_name.append("localhost:10006")
server_name.append("localhost:10007")
server_name.append("localhost:10008")
server_name.append("localhost:10009")
running_server = []
running_pid = []
def create(addr):
    command = "kv-server -addr " + addr
    proc = subprocess.Popen(command.split(" "), shell=False)
    return proc.pid

def kill(pid):
    subprocess.call(["kill", "-15", "%d" % pid])

def start(number):
    global running_server
    to_run = [i for i in server_name if i not in running_server]
    if number > len(to_run):
        number = len(to_run)
    for i in to_run[:number]:
        running_server.append(i)
        running_pid.append(create(i))
def end(number):
    global running_server
    if number > len(running_pid):
        number = len(running_pid)
    to_terminate = running_pid[:number]
    stop_running = running_server[:number]


    for i in to_terminate:
        kill(i)
    running_server = [i for i in running_server if i not in stop_running]


def main():
    while True:
        print "running_server_number:" + str(len(running_server)), \
                                        running_server
        print "$",  
        raw_command = raw_input()
        raw_command = raw_command.split(" ")
        try:
            command = raw_command[0]
            number = int(raw_command[1])
        except:
            continue
        if command == "start":
            start(number)
        if command == "end":
            end(number)
if __name__ == "__main__":
    main()
