#!/bin/bash -x

body='{
  "request": {
    "branch": "phinze/travis",
    "config": {
      "script": [
        "./scripts/travis.sh"
      ],
      "branches": {
        "only": ["phinze/travis"]
      },
      "env": {
        "matrix": [
          "TESTARGS=\"-run TestAccAWSA\"",
          "TESTARGS=\"-run TestAccAWSC\"",
          "TESTARGS=\"-run TestAccAWSD\"",
          "TESTARGS=\"-run TestAccAWSEBS\"",
          "TESTARGS=\"-run TestAccAWSEFS\"",
          "TESTARGS=\"-run TestAccAWSEIP\"",
          "TESTARGS=\"-run TestAccAWSELB\"",
          "TESTARGS=\"-run TestAccAWSENI\"",
          "TESTARGS=\"-run TestAccAWSEcs\"",
          "TESTARGS=\"-run TestAccAWSEla\"",
          "TESTARGS=\"-run TestAccAWSF\"",
          "TESTARGS=\"-run TestAccAWSG\"",
          "TESTARGS=\"-run TestAccAWSI\"",
          "TESTARGS=\"-run TestAccAWSK\"",
          "TESTARGS=\"-run TestAccAWSL\"",
          "TESTARGS=\"-run TestAccAWSM\"",
          "TESTARGS=\"-run TestAccAWSN\"",
          "TESTARGS=\"-run TestAccAWSP\"",
          "TESTARGS=\"-run TestAccAWSR\"",
          "TESTARGS=\"-run TestAccAWSS\"",
          "TESTARGS=\"-run TestAccAWSU\"",
          "TESTARGS=\"-run TestAccAWSV\""
        ],
        "global": [
          "TEST=./builtin/providers/aws",
          "AWS_DEFAULT_REGION=us-west-2",
          {"secure": "Lf3/ppRqqhZpL23lqSt9dvJjbXy6p1ootSgssPXQYlaH/i5S7cKKvyAhW9+VH9fecA3hmc6906dVBYM1LIeb5hqa37S+cIZq9OER8gwCMAGbpwPqJUeoMopZtwArClA8i4dXWudHRAnYO8nMGvEWN/M+V9VrH4jOAU+3mccVUfw="},
          {"secure": "akk1qlzrwvxNi4oobZm6nkKeXcWsdZt68NFX0C1zRe02VZcz6XoanZZ5u2lu2qaxIhxIc/5Nf+ogm2liitbhspdR9vz47ZXFMCyPY0x71Uu5K8PAW8bh611UaIbJoQSEt3+Yz5ZDJDtgLWGRNn82w8rd1WDwwKNyD9HPYGPFjv0="},
          {"secure": "g3Mfj+KMna5tVKY1o14uWqleKAml1WxOVQ0yxTNXP1c5Y0Qz2OP9F4Z0112p/5So7TW8gpsVHWDlH23g881ZjF1sWVWP4Ubjc+y7mXdxkw/YIoZ5ObnDTj3opiiMesWuMA9AeJ+d5hM3OfTIcL/xpLafU9qSp7BhRHs9hvcO29c="}
        ]
      }
    }
  }
}'

echo $body

repo='hashicorp%2Fterraform'

curl -v -X POST \
     -H "Content-Type: application/json" \
     -H "Accept: application/json" \
     -H "Travis-API-Version: 3" \
     -H "Authorization: token ${TRAVIS_TOKEN}" \
     -d "${body}" \
     https://api.travis-ci.org/repo/${repo}/requests
