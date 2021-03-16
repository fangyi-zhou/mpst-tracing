This repos contains preliminary investigations of using opentelemetry and multiparty session types to monitor protocol conformance:

# Layout

- `twobuyer` is an implementation of Two Buyer Protocol in the MPST literature.
  
  Calls to sending and receiving is traced with an opentelemetry tracer.
  Tracer can be changed to export to stdout, open telemetry format, or Jaeger, see `twobuyer/two_buyer.go`
  
- `exporters/mpstconformancemonitoringexporter` is a WIP implementation of a MPST conformance monitor, as an opentelemetry exporter.
  
  To build a collector, follow the instruction on https://github.com/observatorium/opentelemetry-collector-builder and use `manifest.yaml` for the config.
  To run the built collector, use the `config.yaml` for the config.

- `globaltype` is a simple implementation of multiparty session types and their semantics.

- `pedro` is an exploration of petri net semantics of multiparty session types.
