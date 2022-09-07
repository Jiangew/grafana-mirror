import React from 'react';

export default function CheatSheet() {
  return (
    <>
      <h2 id="jaeger-cheat-sheet">Jaeger Cheat Sheet</h2>
      <p>
        Documentation for the Jaeger data source can be found{' '}
        <a
          href="https://grafana.com/docs/grafana/latest/datasources/jaeger"
          target="_blank"
          rel="noreferrer"
          style={{ color: '#1F62E0' }}
        >
          here,
        </a>{' '}
        but here is a quick overview of the query types:
      </p>

      <p>
        Search - filter traces by service name. Addtionally, you can filter by tags or min/max duration, as well as
        limit the number of traces that are returned.
      </p>

      <p>TraceID - if you have a trace ID, simply enter the trace ID to see the trace.</p>

      <p>
        JSON File - you can upload a JSON file that contains a single trace to visualize it. If the file has multiple
        traces then the first trace is used for visualization. Click{' '}
        <a
          href="https://grafana.com/docs/grafana/latest/datasources/jaeger/#upload-json-trace-file"
          target="_blank"
          rel="noreferrer"
          style={{ color: '#1F62E0' }}
        >
          here
        </a>{' '}
        to see an example of a valid JSON file.
      </p>
    </>
  );
}
