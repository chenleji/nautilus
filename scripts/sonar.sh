#!/bin/bash

cat << EOF >../sonar-project.properties
sonar.projectKey=${1}
sonar.projectName=${1}
sonar.projectVersion=${2}

# GoLint report path, default value is report.xml
sonar.golint.reportPath=report.xml

# Cobertura like coverage report path, default value is coverage.xml
sonar.coverage.reportPath=coverage.xml

# if you want disabled the DTD verification for a proxy problem for example, true by default
sonar.coverage.dtdVerification=false

# JUnit like test report, default value is test.xml
sonar.test.reportPath=test.xml
sonar.sources=./
sonar.tests=./
sonar.test.inclusions=**/**_test.go
sonar.sources.inclusions=**/**.go

# sonar server address
sonar.host.url=${3}
EOF

cd ..
sonar-scanner
