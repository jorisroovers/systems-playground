import xml.dom.minidom


from ncclient import manager

host = "localhost"
with manager.connect(host=host, port=830, username="netconf", password="netconf", hostkey_verify=False) as m:
    c = m.get_config(source='running').data_xml
    dom = xml.dom.minidom.parseString(c)
    print dom.toprettyxml()
