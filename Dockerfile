FROM golang

RUN useradd -m oslo

RUN mkdir /build

RUN chown -Rvf oslo: /build

USER oslo


WORKDIR /build

COPY . .

RUN go build

RUN go install

ENTRYPOINT ["oslo"]
