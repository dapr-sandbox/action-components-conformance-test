# Dapr Conformance tests

A github action for running conformance test against a pluggable component.

```yaml
- name: Conformance Tests
  uses: dapr-sandbox/action-components-conformance-test@v0.1.0
  with:
    type: state
    socket: /tmp/socket.sock
    metadata: | ## component-specific init metadata
      timeout: 10s
    operations: | ## If not provided all operations will be tested
      - get
```

A component listening to the specified socket is required, see a complete example below:

```yaml
name: conformance-test

on:
  push:
    branches:
      - main

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - name: Run the component
        shell: bash
        run: docker run -d -it --env DAPR_COMPONENT_SOCKET_PATH=/tmp/socket.sock -v /tmp:/tmp tmacam/dapr-memstore-java:latest

      - name: Conformance Tests
        uses: dapr-sandbox/action-components-conformance-test@v0.1.0
        with:
          type: state
          socket: /tmp/socket.sock
          metadata: |
            timeout: 10s
            concurrency: 1
          operations: |
            - get
            - set
```
