# IAMlist

This application uses data from [trackiam](https://github.com/glassechidna/trackiam) to show a list of all AWS IAM permissions and make them searchable.

This app is deployed to [iamlist.xtyx.nl](https://iamlist.xtyx.nl/). This is redeployed nightly with the newest version of trackiam.

## Development

Since trackiam is included as a submodule, you'll need to initialize submodules. This can be done with `git submodule update --init --recursive` or `npm run build:submodules`.

After that you'll need to generate the json file containing all permissions, this can be done by running `npm run build:iamjson`.

### Basic steps for a dev environment
```
npm i
npm run build:submodules
npm run build:iamjson
npm run dev
```