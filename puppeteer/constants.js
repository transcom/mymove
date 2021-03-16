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
  fileSize: { type: 'string', enum: ['small', 'medium', 'large'] },
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
  average: {
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

const scenarios = {
  too: 'too-orders-document-viewer',
  tio: 'tio-payment-requests-document-viewer',
};

const measurementTypes = {
  totalDuration: 'total-duration',
  networkComparison: 'network-comparison',
  fileDuration: 'file-duration',
};

// All file "tshirt" sizes available
const fileSizes = {
  small: 'small',
  medium: 'medium',
  large: 'large',
};

// All file sizes available
const fileList = {
  small: '150Kb',
  medium: '2mb',
  large: '25mb',
};

const fileSizeList = ['small', 'medium', 'large'];

// File sizes associated to different moves by move code/locator
const fileSizeMoveCodes = {
  small: 'S150KB',
  medium: 'MED2MB',
  large: 'LG25MB',
};

// File sizes associated to different payment requests by id
const fileSizePaymentRequestIds = {
  small: '68034aa3-831c-4d2d-9fd4-b66bc0cc5130',
  medium: '4de88d57-9723-446b-904c-cf8d0a834687',
  large: 'aca5cc9c-c266-4a7d-895d-dc3c9c0d9894',
};

// All speeds available
const speeds = ['fast', 'average', 'slow'];

module.exports = {
  schema,
  networkProfiles,
  measurementTypes,
  fileSizes,
  fileSizeList,
  fileList,
  fileSizeMoveCodes,
  fileSizePaymentRequestIds,
  speeds,
  scenarios,
};
