# Kubernetes Image Mapper

[![REUSE status](https://api.reuse.software/badge/github.com/SAP/image-mapper)](https://api.reuse.software/info/github.com/SAP/image-mapper)

## About this project

This service can act as a [Mutating Kubernetes Admission Webhook](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers) for pods, and allows to dynamically adjust the images used by the containers of a pod, according to configurable rules.

All pods for which the admission webhook is called by the Kubernetes API server are subject to the replacement (if pods should be excluded, this has to be done by selectors in the webhook registration).
For each of the pod's containers, the replacement rules are evaluated, and the first matching rule defines the replacement for the image. Rules are specified in the file provided by command line switch `-mapping-file`. The file should contain a JSON array in the following form:

```
[
  {
    "pattern": "(.+/my-image):.*",
    "replacement": "$1:latest"
  },
  {
    "pattern": "some-registry/.*",
    "replacement": "other-registry/${repository}:${tag}"
  },
  ...
]
```

The pattern can be an arbitrary regular expressions (go syntax) which will be automatically anchored (so adding anchors is not necessary).
If the pattern contains any capturing groups, then the according matches can be used in the replacement as `$1`, `$2`, ..., as usual.
If it does not, then the variables `${registry}`, `${repository}` and `${tag}` will be populated, and can be used in the replacement.

To simplify the rules, the image will be normalized before the rule processing happens, in the following sense: 
- images which do not specify a tag, will be implicitly matched with suffix :latest
- images which do not specify a registry (i.e. Docker hub) will be implicitly matched with prefix docker.io/.

If at least one image was replaced, then configurable labels or annotations can be added, as specified via the command line arguments `-add-label-if-modified` and  `-add-annotation-if-modified` (which can be repeated), in the usual format `key=value`.

Note: in case this webhook has to reliably work with pods that are created or mutated by other webhooks, this one probably has to be registered with `reinvocationPolicy: IfNeeded`.

**Command line flags**

|Flag                         |Optional|Default|Description                                                 |
|-----------------------------|--------|-------|------------------------------------------------------------|
|-bind-address string         |yes     |:2443  |Webhook bind address                                        |
|-tls-key-file                |no      |-      |File containing the TLS private key used for SSL termination|
|-tls-cert-file               |no      |-      |File containing the TLS certificate matching the private key|
|-mapping-file                |no      |-      |File containing the mapping rules                           |
|-add-label-if-modified       |yes     |-      |Label to be set if pod was mutated (can be repeated)        |
|-add-annotation-if-modified  |yes     |-      |Annotation to be set if pod was mutated (can be repeated)   |

## Requirements and Setup

The recommended deployment method is to use the [Helm chart](https://github.com/sap/image-mapper-helm):

```bash
helm upgrade -i image-mapper oci://ghcr.io/sap/image-mapper-helm/image-mapper
```

## Documentation

The API reference is here: [https://pkg.go.dev/github.com/sap/image-mapper](https://pkg.go.dev/github.com/sap/image-mapper).

## Support, Feedback, Contributing

This project is open to feature requests/suggestions, bug reports etc. via [GitHub issues](https://github.com/SAP/image-mapper/issues). Contribution and feedback are encouraged and always welcome. For more information about how to contribute, the project structure, as well as additional contribution information, see our [Contribution Guidelines](CONTRIBUTING.md).

## Code of Conduct

We as members, contributors, and leaders pledge to make participation in our community a harassment-free experience for everyone. By participating in this project, you agree to abide by its [Code of Conduct](https://github.com/SAP/.github/blob/main/CODE_OF_CONDUCT.md) at all times.

## Licensing

Copyright 2025 SAP SE or an SAP affiliate company and image-mapper contributors. Please see our [LICENSE](LICENSE) for copyright and license information. Detailed information including third-party components and their licensing/copyright information is available [via the REUSE tool](https://api.reuse.software/info/github.com/SAP/image-mapper).
