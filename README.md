# gcpcc
[![CodeQL](https://github.com/dacbd/gcpcc/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/dacbd/gcpcc/actions/workflows/codeql-analysis.yml)

Simple tool to print out the current number of compute instances on GCP

Created to to prefrom GitHub Actions powered automatic checks.

Did someone leave a compute instance on that shouldn't be?

## Config
Currently none, authentication credentials are read from the environment.

## Authentication to GCP

There are 3 recomended approaches:
1. Pass GCP credentials JSON directly via `GOOGLE_APPLICATION_CREDENTIALS_DATA`
1. The ["typical"](https://github.com/google-github-actions/auth#authenticating-via-service-account-key-json-1) approach of `GOOGLE_APPLICATION_CREDENTIALS` (path to JSON file)
1. Using [GitHub Actions OIDC](https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/configuring-openid-connect-in-google-cloud-platform) with GCP [see here](https://github.com/google-github-actions/auth#authenticating-via-workload-identity-federation-1)

#### Method 1
```yml
# ...
steps:
  - uses: dacbd/gcpcc@v1
    env:
      GOOGLE_APPLICATION_CREDENTIALS_DATA: ${{ secrets.GCP_SA_KEY_JSON }}
```

#### Method 2
```yml
# ...
steps:
  - uses: google-github-actions/auth@v0
    with:
      credentials_json: ${{ secrets.gcp_sa_key_json }}
  - uses: dacbd/gcpcc@v1
```

#### Method 3
```yml
# ...
# Add "id-token" with the intended permissions.
permissions:
  contents: 'read'
  id-token: 'write'
steps:
  - uses: google-github-actions/auth@v0
    with:
      workload_identity_provider: 'projects/123456789/locations/global/workloadIdentityPools/my-pool/providers/my-provider'
      service_account: 'my-service-account@my-project.iam.gserviceaccount.com'
  - uses: dacbd/gcpcc@v1
```

## Outputs
| outputs | value |
| ------- | ----- |
| total   | `int` - total number of compute instance |

## Usage
Basic example:
```yml
name: Check ML Training instances
on:
  schedule:
    cron:
     - ''
jobs:
  check:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3 # not technically required
    - uses: google-github-actions/auth@v0
      with:
        credentials_json: ${{ secrets.GCP_SA_KEY_JSON }}
    - uses: dacbd/gcpcc@v1
      id: gcpcc
    - uses: dacbd/create-issue-action@v1
      if: steps.gcpcc.outputs.total != 0
      with:
        token: ${{ github.token }}
        title: Instance left on in `${{ env.GCP_PROJECT }}`
        assignees: dacbd,some_github_username
        body: |
          Automatic check found `${{ steps.gcpcc.outputs.total }}` instance\s left on.
```
## Permissions
TODO
