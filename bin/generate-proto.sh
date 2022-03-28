#!/usr/bin/env bash

set -euo pipefail

cd api/proto
buf generate
cd ../..