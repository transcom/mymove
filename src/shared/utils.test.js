import * as utils from './utils';

describe('utils', () => {
  describe('upsert', () => {
    const item = { id: 'foo', name: 'something' };
    describe('when upserting a new item to an array', () => {
      const arr = [
        { id: 'bar', name: 'foo' },
        { id: 'baz', name: 'baz' },
      ];
      utils.upsert(arr, item);
      it('should be appended to the array', () => {
        expect(arr).toEqual([{ id: 'bar', name: 'foo' }, { id: 'baz', name: 'baz' }, item]);
      });
    });
    describe('when upserting an update to an array', () => {
      const arr = [
        { id: 'foo', name: 'foo' },
        { id: 'baz', name: 'baz' },
      ];
      utils.upsert(arr, item);
      it('should be appended to the array', () => {
        expect(arr).toEqual([
          { id: 'foo', name: 'something' },
          { id: 'baz', name: 'baz' },
        ]);
      });
    });
  });

  describe('fetch Active', () => {
    describe('when there are no foos', () => {
      const foos = null;
      const res = utils.fetchActive(foos);
      it('should return null', () => {
        expect(res).toBeNull();
      });
    });
    describe('when there are some active and some inactive foos', () => {
      const foos = [
        { status: 'CANCELED', id: 'foo0' },
        { status: 'DRAFT', id: 'foo1' },
        { status: 'SUBMITTED', id: 'foo2' },
      ];
      const res = utils.fetchActive(foos);
      it('should return the first active foo', () => {
        expect(res.id).toEqual('foo1');
      });
    });
    describe('when there are only inactive foos', () => {
      const foos = [
        { status: 'CANCELED', id: 'foo' },
        { status: 'COMPLETED', id: 'foo' },
      ];
      const res = utils.fetchActive(foos);
      it('should return null', () => {
        expect(res).toBeNull();
      });
    });
  });

  describe('fetch Active Shipment', () => {
    describe('when there are no foos', () => {
      const foos = null;
      const res = utils.fetchActiveShipment(foos);
      it('should return null', () => {
        expect(res).toBeNull();
      });
    });
    describe('when there are some active and some inactive foos', () => {
      const foos = [
        { status: 'CANCELED', id: 'foo0' },
        { status: 'DRAFT', id: 'foo1' },
        { status: 'SUBMITTED', id: 'foo2' },
      ];
      const res = utils.fetchActiveShipment(foos);
      it('should return the first active foo', () => {
        expect(res.id).toEqual('foo1');
      });
    });
    describe('when there are only inactive foos', () => {
      const foos = [{ status: 'CANCELED', id: 'foo' }];
      const res = utils.fetchActiveShipment(foos);
      it('should return null', () => {
        expect(res).toBeNull();
      });
    });
  });

  it('check if 2nd and 3rd addresses should be cleared from prime shipment create payload', () => {
    const ppmValues = {
      shipmentType: 'PPM',
      ppmShipment: {
        hasSecondaryPickupAddress: 'false',
        hasTertiaryPickupAddress: 'false',
        hasSecondaryDestinationAddress: 'false',
        hasTertiaryDestinationAddress: 'false',
        secondaryPickupAddress: '',
        tertiaryPickupAddress: '',
        secondaryDestinationAddress: '',
        tertiaryDestinationAddress: '',
      },
    };
    const hhgValues = {
      shipmentType: 'HHG',
      hasSecondaryPickupAddress: 'false',
      hasTertiaryPickupAddress: 'false',
      hasSecondaryDestinationAddress: 'false',
      hasTertiaryDestinationAddress: 'false',
      secondaryPickupAddress: '',
      tertiaryPickupAddress: '',
      secondaryDestinationAddress: '',
      tertiaryDestinationAddress: '',
    };

    const updatedPPMValues = utils.checkAddressTogglesToClearAddresses(ppmValues);
    expect(updatedPPMValues).toEqual({
      shipmentType: 'PPM',
      ppmShipment: {
        hasSecondaryPickupAddress: 'false',
        hasTertiaryPickupAddress: 'false',
        hasSecondaryDestinationAddress: 'false',
        hasTertiaryDestinationAddress: 'false',
        secondaryPickupAddress: {},
        tertiaryPickupAddress: {},
        secondaryDestinationAddress: {},
        tertiaryDestinationAddress: {},
      },
    });

    const updatedHHGValues = utils.checkAddressTogglesToClearAddresses(hhgValues);
    expect(updatedHHGValues).toEqual({
      shipmentType: 'HHG',
      hasSecondaryPickupAddress: 'false',
      hasTertiaryPickupAddress: 'false',
      hasSecondaryDestinationAddress: 'false',
      hasTertiaryDestinationAddress: 'false',
      secondaryPickupAddress: {},
      tertiaryPickupAddress: {},
      secondaryDestinationAddress: {},
      tertiaryDestinationAddress: {},
    });
  });
});

describe('appendTimestampToFilename', () => {
  // Define the fileTypeMap for reference in tests
  const fileTypeMap = {
    'application/pdf': 'pdf',
    'image/png': 'png',
    'image/jpeg': 'jpg',
    'image/jpg': 'jpg',
    'image/gif': 'gif',
  };
  let mockDate;
  let expectedTimestamp;
  beforeEach(() => {
    // Generate a dynamic mock date (e.g., current time or random past date)
    mockDate = new Date();
    jest.spyOn(global, 'Date').mockImplementation(() => mockDate);
    // Derive expected timestamp dynamically based on the function's logic
    expectedTimestamp = mockDate
      .toISOString()
      .replace(/[-T:.Z]/g, '')
      .slice(0, 14);
  });
  afterEach(() => {
    jest.restoreAllMocks();
  });
  it('appends timestamp to filename without extension', () => {
    const file = new File([''], 'test-file', { type: 'application/octet-stream' });
    const result = utils.appendTimestampToFilename(file);
    expect(result.name).toBe(`test-file-${expectedTimestamp}`);
    expect(result.type).toBe('application/octet-stream');
    expect(result instanceof File).toBe(true);
  });
  it('handles filenames with multiple dots', () => {
    const file = new File([''], 'test.file.with.dots.jpg', { type: 'image/jpeg' });
    const result = utils.appendTimestampToFilename(file);
    expect(result.name).toBe(`test.file.with.dots-${expectedTimestamp}.jpg`);
    expect(result.type).toBe('image/jpeg');
    expect(result instanceof File).toBe(true);
  });
  it('handles whitespace filename', () => {
    const file = new File([''], ' .pdf', { type: 'image/pdf' });
    const result = utils.appendTimestampToFilename(file);
    expect(result.name).toBe(` -${expectedTimestamp}.pdf`);
    expect(result.type).toBe('image/pdf');
    expect(result instanceof File).toBe(true);
  });
  it('empty filename with extension', () => {
    const file = new File([''], '.pdf', { type: 'image/pdf' });
    const result = utils.appendTimestampToFilename(file);
    expect(result.name).toBe(`-${expectedTimestamp}.pdf`);
    expect(result.type).toBe('image/pdf');
    expect(result instanceof File).toBe(true);
  });
  it('empty filename and no extension', () => {
    const file = new File([''], '', { type: 'image/pdf' });
    const result = utils.appendTimestampToFilename(file);
    expect(result.name).toBe(`-${expectedTimestamp}`);
    expect(result.type).toBe('image/pdf');
    expect(result instanceof File).toBe(true);
  });
  it('filename no extension', () => {
    const file = new File([''], 'helloworld', { type: 'image/pdf' });
    const result = utils.appendTimestampToFilename(file);
    expect(result.name).toBe(`helloworld-${expectedTimestamp}`);
    expect(result.type).toBe('image/pdf');
    expect(result instanceof File).toBe(true);
  });
  // Test each file type in fileTypeMap
  Object.entries(fileTypeMap).forEach(([mimeType, extension]) => {
    it(`appends timestamp to filename with ${extension} extension`, () => {
      const file = new File([''], `test-file.${extension}`, { type: mimeType });
      const result = utils.appendTimestampToFilename(file);
      expect(result.name).toBe(`test-file-${expectedTimestamp}.${extension}`);
      expect(result.type).toBe(mimeType);
      expect(result instanceof File).toBe(true);
    });
  });
});
