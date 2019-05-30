
# Copyright 2018-present Open Networking Foundation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# docker build -t opencord/sadis-server:candidate .

FROM golang:1.10-stretch as builder
MAINTAINER Open Networking Foundation <info@opennetworking.org>

WORKDIR /go
ADD . /go/src/gerrit.opencord.org/sadis-server
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/entry-point gerrit.opencord.org/sadis-server

# The amd64 and arm64 are supported by alpine:3.6.
FROM alpine:3.9
MAINTAINER Open Networking Foundation <info@opennetworking.org>

COPY --from=builder /build/entry-point /service/entry-point

EXPOSE 8000

WORKDIR /service
ENTRYPOINT ["/service/entry-point"]
