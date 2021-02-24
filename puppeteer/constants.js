const schema = {
  launch: {
    type: 'object',
    properties: {
      headless: { type: 'boolean' },
    },
  },
  device: { type: 'string' },
  emulate: {
    type: 'object',
    properties: {
      viewport: {
        type: 'object',
        properties: {
          width: { type: 'number' },
          height: { type: 'number' },
          deviceScaleFactor: { type: 'number' },
          isMobile: { type: 'boolean' },
          hasTouch: { type: 'boolean' },
          isLandscape: { type: 'boolean' },
        },
      },
      userAgent: { type: 'string' },
    },
  },
  network: { type: 'string', enum: ['fast', 'medium', 'slow'] },
  throttling: {
    type: 'object',
    properties: {
      offline: { type: 'boolean' },
      download: { type: 'number' },
      upload: { type: 'number' },
      latency: { type: 'number' },
    },
  },
  cpuSlowdownRate: { type: 'number', default: 1 },
};

// Chrome devtools expects throughput in bytes https://github.com/GoogleChrome/lighthouse/blob/master/lighthouse-core/lib/emulation.js#L70
// Puppeteer only has 2 network presets https://github.com/puppeteer/puppeteer/blob/main/src/common/NetworkConditions.ts#L27
const networkProfiles = {
  fast: {
    offline: false,
    download: (10 * 1024 * 1024) / 8,
    upload: (10 * 1024 * 1024) / 8,
    latency: 20,
  },
  medium: {
    offline: false,
    download: (5 * 1024 * 1024) / 8,
    upload: (5 * 1024 * 1024) / 8,
    latency: 150,
  },
  slow: {
    offline: false,
    download: (1024 * 1024) / 8,
    upload: (1024 * 1024) / 8,
    latency: 250,
  },
};

const measurementTypes = {
  totalDuration: 'total-duration',
  networkComparison: 'network-comparison',
};

// All file sizes available
const fileSizes = {
  small: 'small',
  medium: 'medium',
  large: 'large',
};

// File sizes associated to different moves by move code/locator
const fileSizeMoveCodes = {
  small: 'S150KB',
  medium: 'MED2MB',
  large: 'LG25MB',
};

// All speeds available
const speeds = ['fast', 'medium', 'slow'];

module.exports = { schema, networkProfiles, measurementTypes, fileSizes, fileSizeMoveCodes, speeds };
