// Thanks sir: https://stackoverflow.com/questions/33524537/good-zookeeper-hello-world-program-with-java-client/33603040
package com.mycompany.zkexample;

import java.util.Date;
import java.util.List;
import java.util.concurrent.CountDownLatch;

import org.apache.zookeeper.CreateMode;
import org.apache.zookeeper.WatchedEvent;
import org.apache.zookeeper.Watcher;
import org.apache.zookeeper.Watcher.Event.KeeperState;
import org.apache.zookeeper.ZooDefs.Ids;
import org.apache.zookeeper.ZooKeeper;

public class ZkConnect {
    private ZooKeeper zk;
    private CountDownLatch connSignal = new CountDownLatch(0);

    public ZooKeeper connect(String host) throws Exception {
        zk = new ZooKeeper(host, 5000 /* timeout */, new Watcher() {
            public void process(WatchedEvent event) {
                if (event.getState() == KeeperState.SyncConnected) {
                    connSignal.countDown();
                }
            }
        });
        connSignal.await();
        return zk;
    }

    public void close() throws InterruptedException {
        zk.close();
    }

    public void createNode(String path, byte[] data) throws Exception {
        zk.create(path, data, Ids.OPEN_ACL_UNSAFE, CreateMode.PERSISTENT);
    }

    public void updateNode(String path, byte[] data) throws Exception {
        zk.setData(path, data, zk.exists(path, true).getVersion());
    }

    public void deleteNode(String path) throws Exception {
        zk.delete(path,  zk.exists(path, true).getVersion());
    }

    public static void main (String args[]) throws Exception {
        ZkConnect connector = new ZkConnect();
        ZooKeeper zk = connector.connect("zookeeper" /* hostname */);

        // Create new node
        String newNode = "/myDate" + new Date();
        connector.createNode(newNode, new Date().toString().getBytes());

        // List existing nodes
        List<String> zNodes = zk.getChildren("/", true);
        for (String zNode: zNodes)
        {
           System.out.println("ChildrenNode " + zNode);
        }

        System.out.println("**************************************************");
        // Get Node
        byte[] data = zk.getData(newNode, true, zk.exists(newNode, true));
        System.out.print("GetData before setting: ");
        for (byte dataPoint: data) {
            System.out.print((char)dataPoint);
        }
        System.out.println();

        // Set Node (=update existing node)
        System.out.print("GetData after setting: ");
        connector.updateNode(newNode, "Modified data".getBytes());
        data = zk.getData(newNode, true, zk.exists(newNode, true));
        for (byte dataPoint: data) {
            System.out.print((char)dataPoint);
        }
        System.out.println();

        // Delete Node
        connector.deleteNode(newNode);
    }

}