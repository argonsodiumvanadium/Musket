#!/usr/bin/bash
export PATH=$PATH:~/Documents/VIPER-master/src/linux
if [[ $1 == "" || $1 == "--help" || $1 == "help" ]]
then
	printf "                    \u001B[1mVIPER_LANG\n"
	printf "                ELIGIBLE ARGUMENTS\n\n\n"
	printf "                     commands\u001B[0m\n"
	echo "run [filename]                runs the specified file"
	echo "run -d          runs default file [lethalityTest.vpr]"
	echo "shell                               opens viper shell"
	printf "\n\n\n"
	printf "                      \u001B[1mflags\u001B[0m\n"
	echo "-r [filename]                runs the specified file"
	echo "-r -d          runs default file [lethalityTest.vpr]"
	echo "-s                                 opens viper shell"
elif [[ $1 == "shell" || $1 == "-s" ]];
then
	./main
elif [[ $1 == "-r" ]]; 
then
	./main run $2	
else
	./main $1 $2
fi
