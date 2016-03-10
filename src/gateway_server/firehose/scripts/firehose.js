envelopeProto = dcodeIO.ProtoBuf.loadProtoFile('dropsonde-protocol/events/envelope.proto').build('events.Envelope')

var ws = new WebSocket('ws://doppler.bosh-east.high.am/firehose/1234');
ws.binaryType = "arraybuffer";

ws.onmessage = function (event) {
  var envelope = envelopeProto.decode(event.data);
  if (envelope.origin == 'rep') {
    console.log(envelope);
  }
};
