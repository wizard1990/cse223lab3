import subprocess
import time

global server_name
global running_server
server_name = [
        "localhost:18067",
        "localhost:18068",
        "localhost:18069",
        "localhost:18070",
        "localhost:18071",
        "localhost:18072",
        "localhost:18073",
        "localhost:18074",
        "localhost:18075",
        "localhost:18076",
        "localhost:18077",
        "localhost:18078",
        "localhost:18079",
        "localhost:18080",
        "localhost:18081",
        "localhost:18082",
        "localhost:18083",
        "localhost:18084",
        "localhost:18085",
        "localhost:18086",
        "localhost:18087",
        "localhost:18088",
        "localhost:18089",
        "localhost:18090",
        "localhost:18091",
        "localhost:18092",
        "localhost:18093",
        "localhost:18094",
        "localhost:18095",
        "localhost:18096",
        "localhost:18097",
        "localhost:18098",
        "localhost:18099",
        "localhost:18100",
        "localhost:18101",
        "localhost:18102",
        "localhost:18103",
        "localhost:18104",
        "localhost:18105",
        "localhost:18106",
        "localhost:18107",
        "localhost:18108",
        "localhost:18109",
        "localhost:18110",
        "localhost:18111",
        "localhost:18112",
        "localhost:18113",
        "localhost:18114",
        "localhost:18115",
        "localhost:18116",
        "localhost:18117",
        "localhost:18118",
        "localhost:18119",
        "localhost:18120",
        "localhost:18121",
        "localhost:18122",
        "localhost:18123",
        "localhost:18124",
        "localhost:18125",
        "localhost:18126",
        "localhost:18127",
        "localhost:18128",
        "localhost:18129",
        "localhost:18130",
        "localhost:18131",
        "localhost:18132",
        "localhost:18133",
        "localhost:18134",
        "localhost:18135",
        "localhost:18136",
        "localhost:18137",
        "localhost:18138",
        "localhost:18139",
        "localhost:18140",
        "localhost:18141",
        "localhost:18142",
        "localhost:18143",
        "localhost:18144",
        "localhost:18145",
        "localhost:18146",
        "localhost:18147",
        "localhost:18148",
        "localhost:18149",
        "localhost:18150",
        "localhost:18151",
        "localhost:18152",
        "localhost:18153",
        "localhost:18154",
        "localhost:18155",
        "localhost:18156",
        "localhost:18157",
        "localhost:18158",
        "localhost:18159",
        "localhost:18160",
        "localhost:18161",
        "localhost:18162",
        "localhost:18163",
        "localhost:18164",
        "localhost:18165",
        "localhost:18166",
        "localhost:18167",
        "localhost:18168",
        "localhost:18169",
        "localhost:18170",
        "localhost:18171",
        "localhost:18172",
        "localhost:18173",
        "localhost:18174",
        "localhost:18175",
        "localhost:18176",
        "localhost:18177",
        "localhost:18178",
        "localhost:18179",
        "localhost:18180",
        "localhost:18181",
        "localhost:18182",
        "localhost:18183",
        "localhost:18184",
        "localhost:18185",
        "localhost:18186",
        "localhost:18187",
        "localhost:18188",
        "localhost:18189",
        "localhost:18190",
        "localhost:18191",
        "localhost:18192",
        "localhost:18193",
        "localhost:18194",
        "localhost:18195",
        "localhost:18196",
        "localhost:18197",
        "localhost:18198",
        "localhost:18199",
        "localhost:18200",
        "localhost:18201",
        "localhost:18202",
        "localhost:18203",
        "localhost:18204",
        "localhost:18205",
        "localhost:18206",
        "localhost:18207",
        "localhost:18208",
        "localhost:18209",
        "localhost:18210",
        "localhost:18211",
        "localhost:18212",
        "localhost:18213",
        "localhost:18214",
        "localhost:18215",
        "localhost:18216",
        "localhost:18217",
        "localhost:18218",
        "localhost:18219",
        "localhost:18220",
        "localhost:18221",
        "localhost:18222",
        "localhost:18223",
        "localhost:18224",
        "localhost:18225",
        "localhost:18226",
        "localhost:18227",
        "localhost:18228",
        "localhost:18229",
        "localhost:18230",
        "localhost:18231",
        "localhost:18232",
        "localhost:18233",
        "localhost:18234",
        "localhost:18235",
        "localhost:18236",
        "localhost:18237",
        "localhost:18238",
        "localhost:18239",
        "localhost:18240",
        "localhost:18241",
        "localhost:18242",
        "localhost:18243",
        "localhost:18244",
        "localhost:18245",
        "localhost:18246",
        "localhost:18247",
        "localhost:18248",
        "localhost:18249",
        "localhost:18250",
        "localhost:18251",
        "localhost:18252",
        "localhost:18253",
        "localhost:18254",
        "localhost:18255",
        "localhost:18256",
        "localhost:18257",
        "localhost:18258",
        "localhost:18259",
        "localhost:18260",
        "localhost:18261",
        "localhost:18262",
        "localhost:18263",
        "localhost:18264",
        "localhost:18265",
        "localhost:18266"]
running_server = []
running_pid = []
def create(addr):
    command = "kv-server -addr " + addr
    proc = subprocess.Popen(command.split(" "), shell=False)
    return proc

def kill(proc):
    #subprocess.call(["kill", "-15", "%d" % pid])
    proc.kill()

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
