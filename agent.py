import json
import sys
from pysnmp.hlapi import *
from pysnmp.entity import engine, config
from pysnmp.carrier.asyncore.dgram import udp
from pysnmp.entity.rfc3413 import cmdrsp, context
from pysnmp.smi import builder, view, compiler, rfc1902

def load_metrics(file_path):
    try:
        with open(file_path, 'r') as f:
            data = json.load(f)
        print(f"Metrics loaded from {file_path}")
        return data
    except Exception as e:
        print(f"Failed to load metrics: {e}")
        sys.exit(1)

def create_snmp_engine(metrics, port):
    try:
        snmpEngine = engine.SnmpEngine()

        config.addSocketTransport(
            snmpEngine,
            udp.domainName,
            udp.UdpTransport().openServerMode(('0.0.0.0', port))
        )

        config.addV1System(snmpEngine, 'my-area', 'public')
        config.addVacmUser(snmpEngine, 1, 'my-area', 'noAuthNoPriv', (1, 3, 6, 1, 2, 1))

        snmpContext = context.SnmpContext(snmpEngine)

        mibBuilder = snmpContext.getMibInstrum().getMibBuilder()
        mibViewController = view.MibViewController(mibBuilder)
        compiler.addMibCompiler(mibBuilder, sources=['http://mibs.snmplabs.com/asn1/@mib@'])

        mibBuilder.loadModules('SNMPv2-MIB')

        def handle_get(oid, *args):
            oid_str = str(oid)
            if oid_str in metrics:
                return rfc1902.OctetString(metrics[oid_str])
            return rfc1902.NoSuchInstance()

        class MyMibInstrumController(cmdrsp.GetCommandResponder):

            def handleMgmtReq(self, snmpEngine, stateReference, contextName, varBinds, *context):
                result = []
                for name, val in varBinds:
                    result.append((name, handle_get(name)))
                self.sendVarBinds(snmpEngine, stateReference, 0, contextName, result)

        MyMibInstrumController(snmpEngine, snmpContext)

        print(f"SNMP Engine created and listening on port {port}")
        return snmpEngine
    except Exception as e:
        print(f"Failed to create SNMP engine: {e}")
        sys.exit(1)

def main():
    if len(sys.argv) != 3:
        print("Usage: agent.py <path_to_json_file> <port>")
        sys.exit(1)

    file_path = sys.argv[1]
    port = int(sys.argv[2])

    print(f"Starting SNMP agent with metrics from {file_path} on port {port}")

    metrics = load_metrics(file_path)
    snmpEngine = create_snmp_engine(metrics, port)

    print(f"SNMP Agent is running on port {port}...")
    snmpEngine.transportDispatcher.jobStarted(1)

    try:
        snmpEngine.transportDispatcher.runDispatcher()
    except Exception as e:
        print(f"Error in SNMP dispatcher: {e}")
        snmpEngine.transportDispatcher.closeDispatcher()
        raise

main()
