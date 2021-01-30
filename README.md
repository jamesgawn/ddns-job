# DDNS Job

A very lightweight Go executable designed to automate updating a DNS record in Route 53 with your current external IP address.

## How to use

The simplest way to use this project is to use the pre-built docker container with a command such as:

```bash
docker run -ti --env AWS_REGION=eu-west-2 --env AWS_ACCESS_KEY_ID=BLAH --env AWS_SECRET_ACCESS_KEY=BLAHKEY --rm ghcr.io/jamesgawn/ddns-job:latest "-z home.gawn.uk" "-d home.gawn.uk" 
```

_There are a variety of means to supply credentials to the environment, see the [AWS CLI configuration and file settings](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html) page for details._ 

In order for the command to work the supplied AWS credentials will need to have the following permissions:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "route53:ChangeResourceRecordSets",
            "Resource": "arn:aws:route53:::hostedzone/*"
        },
        {
            "Effect": "Allow",
            "Action": "route53:ListHostedZones",
            "Resource": "*"
        }
    ]
}
```

Or preferable this more targeted configuration:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "route53:ChangeResourceRecordSets",
            "Resource": "arn:aws:route53:::hostedzone/INSERTZONEID"
        },
        {
            "Effect": "Allow",
            "Action": "route53:ListHostedZones",
            "Resource": "*"
        }
    ]
}
```

_You can find your zone ID by looking in the Route53 console, opening the relevant zone, and expanding the "Hosted zone details" panel._

## How to build

```bash
make build
```

It will be built into a binary here: .bin/ddns-job