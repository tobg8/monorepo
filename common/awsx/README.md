# awsx

AWSx is an auth library built on top of [aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2).

It exposes a Config structure allowing to easily connect to any AWS services (s3, sqs, ...).

## Usage

This lib is ***not to be used alone***. It's a configuration builder to allow usage of
[sqsx](https://github.mpi-internal.com/leboncoin/go/tree/master/common/messaging/sqsx),
[snsx](https://github.mpi-internal.com/leboncoin/go/tree/master/common/messaging/snsx),
[s3x](https://github.mpi-internal.com/leboncoin/go/tree/master/common/filestore/s3x).

Most of the configuration is done on [terraform](https://backstage.mpi-internal.com/docs/lbc/system/kubernetes/pages/how-to/kaas/associate-iam-app/).
The pod has to be configured with the correct role. You can find an example [here](https://review.leboncoin.ci/c/core/infra/k8s-service-accounts/+/174355).

This will inject the correct [environmental variables](https://github.com/aws/aws-sdk-go-v2/blob/main/config/env_config.go#L23) in the pod.

### From leboncoin AWS account

This is the out-of-the-box case. You can specify a specific logger ```WithLogger``` (compatible with logging.Logger)
and httpClient ```WithHttpClient``` (compatible with the standard http.Doer) with builder functions.

### From another AWS account

For this case, you have to define:

- a role on your account who will have access to the feature you need (s3, sns, ...)
- a role on lbc account who will have the right to assume the first role

You can then use the ```WithAssumeRole``` to assume the correct role.
