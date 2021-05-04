# JMX
Simple JMX Experiment, based on https://www.baeldung.com/java-management-extensions

Compiled and tested using OpenJDK JRE 16

```sh
# Compile java code
javac jmxexperiment/*.java

# Run code (will block)
java jmxexperiment.JMXExperiment

# Open JConsole (in separate tab)
jconsole

# You can now select the local process 'jmxexperiment.JMXExperiment' and then find the managed bean
# You can do then do read/write operations to it
```