/* 	acquire key events from input device ("key logger")
	and write it (slightly modified) to fifo
	where a partner program may fetch which key has been pressed
	useful for headless mode, i.e. autostarted service must react to device input without
	other than on its own tty
	
	run with, e.g. 
	sudo ./key-event /dev/input/event5	
	
	modified after
	http://www.thelinuxdaily.com/2010/05/grab-raw-keyboard-input-from-event-device-node-devinputevent/
*/

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <errno.h>
#include <fcntl.h>
#include <dirent.h>
#include <linux/input.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <sys/select.h>
#include <sys/time.h>
#include <termios.h>
#include <signal.h>


// TODO named pipe's path via argv[2]


int main (int argc, char *argv[]) {

	struct input_event ev[64];
	int fd, rd, value, size = sizeof (struct input_event);
	char *device = NULL;
	int pipe;
	char fifo[256] = "/tmp/gofifo";
	char key[4];

	if (argv[1] == NULL) {
		printf("argv[1] must be path to the dev event interface device\n");
		exit (0);
	}

	if ((getuid()) != 0) {
		printf ("You are not root! This may not work...\n");
	}

	if (argc > 1) {
		device = argv[1];
	}

	// open device
	if ((fd = open(device, O_RDONLY)) == -1) {
		printf("%s is not a vaild device.\n", device);
		exit(1);
	}

	/*	use for debugging:
	
		char name[256] = "Unknown";
		ioctl(fd, EVIOCGNAME (sizeof (name)), name);
		printf("Reading From : %s (%s)\n", device, name);
	*/

	while (1) {
	
		// read from device file descriptor
		if ((rd = read(fd, ev, size*64)) < size) {
			exit(3);      
		}

	  	value = ev[0].value;

		// only read the key press event
	  	if (value != ' ' && ev[1].value == 1 && ev[1].type == 1) {
	   		// printf("code [%d]\n", (ev[1].code));
	   		
	   		// convert key code to string "<code>" 
	   		snprintf(key, 4, "%d", (ev[1].code));
	   		// printf("key is now %s\n", key);
	   		
	   	        // write string to named pipe (fifo)
			pipe = open(fifo, O_WRONLY);
			write(pipe, key, 4);
			close(pipe);
	  	}
	}

	return 0;
} 
