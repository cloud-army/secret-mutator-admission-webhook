#!bin/sh
# WARNING: This is for testing purposes only, 
# never use this for production
printenv

sleep 50s

# INFO: Place here the main execution of your application:
# Some examples:
# java -jar app.jar
# python app.py
# go run main.go
# The application (app) will finally receive the injection environment variables
# The variables will be available only to the thread running the application, 
# not as kubernetes secrets, and not printable in the pod
