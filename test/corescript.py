#!/usr/bin/python

from core import pycore

session = pycore.Session(persistent=True)
node1 = session.addobj(cls=pycore.nodes.CoreNode, name="n1")
node2 = session.addobj(cls=pycore.nodes.CoreNode, name="n2")
hub1 = session.addobj(cls=pycore.nodes.HubNode, name="hub1")
node1.newnetif(hub1, ["10.0.0.1/24"])
node2.newnetif(hub1, ["10.0.0.2/24"])

node1.nodefilecopy("scrooge", "../scrooge")

node1.icmd(["./scrooge"])
node1.icmd(["ping", "-c", "5", "10.0.0.2"])
session.shutdown()
