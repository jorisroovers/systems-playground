def test_passwd_file(host):
    myfile = host.file("/tmp/foo.txt")
    assert myfile.contains("This is a template!")
