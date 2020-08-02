# IXXAT VCI3 Golang package

[github.com/amdf/ixxatvci3](https://github.com/amdf/ixxatvci3) is a CAN bus golang package supporting IXXAT VCI3 interface for Windows.

`ixxatvci3` works with IXXAT USB-to-CAN devices.

## Installation

```bash
go get github.com/amdf/ixxatvci3
```

_ixxatvci3_ is *cgo* package.
If you want to build your app using ixxatvci3, you need gcc. [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) is recommended for Windows.

***Important: because this is a `CGO` enabled package you are required to set the environment variable `CGO_ENABLED=1` and have a `gcc` compile present within your path.***

## Examples
See https://github.com/amdf/ixxatvci3-examples
