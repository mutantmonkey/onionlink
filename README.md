# onionlink

**NOTE:** This app has not been updated since 2015 and should be considered
deprecated. If you need this functionality, I recommend that you use
[OnionShare](https://github.com/micahflee/onionshare).

Very simple peer-to-peer file sharing using Onion services.

## Prerequisites
* Tor 0.2.7+ with a ControlSocket set to /run/tor/control and a corresponding
  auth cookie at /run/tor/control.authcookie

## Usage
1. `go get github.com/mutantmonkey/onionlink`
2. `onionlink file1 [file2] [...]`
3. Share the provided links.
