FROM trussworks/circleci-docker-primary:9b7bfbf6bfae544566ade0c896f8aacbd16281e5

ENV GOPATH=/home/circleci/go
ENV GOBIN=$GOPATH/bin
ENV PATH=$GOBIN:$PATH

COPY --chown=circleci:circleci . /home/circleci/milmove
WORKDIR /home/circleci/milmove
RUN mkdir -p /home/circleci/.cache/
RUN chown circleci:circleci /home/circleci/.cache/

RUN make go_deps_update
RUN make client_deps_update
RUN git --no-pager status && git --no-pager diff --ignore-all-space --color
