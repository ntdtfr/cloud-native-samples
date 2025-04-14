// babel.config.js
module.exports = {
  presets: [
    [
      '@babel/preset-env',
      {
        targets: {
          node: 'current',
        },
        // Use polyfills for modern JS features when needed
        useBuiltIns: 'usage',
        corejs: 3,
      },
    ],
  ],
  plugins: [
    // Transform class properties
    '@babel/plugin-proposal-class-properties',

    // Optional chaining (?.) and nullish coalescing operators (??)
    '@babel/plugin-proposal-optional-chaining',
    '@babel/plugin-proposal-nullish-coalescing-operator',

    // Supports dynamic import() syntax
    '@babel/plugin-syntax-dynamic-import',

    // Transform async/await to backwards compatible code
    '@babel/plugin-transform-runtime',
  ],
  // Ignore certain files during transformation
  ignore: [
    'node_modules',
    'dist',
  ],
  // Source maps for better debugging
  sourceMaps: 'inline',
  // Cache compilation results for better performance
  cacheDirectory: true,
};
