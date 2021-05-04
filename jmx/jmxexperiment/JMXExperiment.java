package jmxexperiment;


import javax.management.InstanceAlreadyExistsException;
import javax.management.MBeanServer;
import javax.management.MalformedObjectNameException;
import javax.management.NotCompliantMBeanException;
import javax.management.ObjectName;
import java.lang.management.ManagementFactory;

public class JMXExperiment {

    public static void main(String[] args) {

        System.out.println("hello world!");

        try {
            // The object Name doesn't need to follow the package qualified name, it can be completely different
            // THe part after the colon is just a comma-separated listed of key-value pairs
            ObjectName objectName = new ObjectName("com.jorisroovers.jmx.experiment:type=basic,name=game");
            MBeanServer server = ManagementFactory.getPlatformMBeanServer();
            server.registerMBean(new Game(), objectName);
        } catch (Exception e) {
            e.printStackTrace();
        }

        System.out.println("Registration for Game mbean with the platform server is successfull");
        System.out.println("Please open jConsole to access Game mbean");

        while (true) {
            // to ensure application does not terminate
        }
    }

}