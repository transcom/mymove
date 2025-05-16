import appendTimestampToFilename from 'utils/fileUpload';

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
    const result = appendTimestampToFilename(file);
    expect(result.name).toBe(`test-file-${expectedTimestamp}`);
    expect(result.type).toBe('application/octet-stream');
    expect(result instanceof File).toBe(true);
  });
  it('handles filenames with multiple dots', () => {
    const file = new File([''], 'test.file.with.dots.jpg', { type: 'image/jpeg' });
    const result = appendTimestampToFilename(file);
    expect(result.name).toBe(`test.file.with.dots-${expectedTimestamp}.jpg`);
    expect(result.type).toBe('image/jpeg');
    expect(result instanceof File).toBe(true);
  });
  it('handles whitespace filename', () => {
    const file = new File([''], ' .pdf', { type: 'image/pdf' });
    const result = appendTimestampToFilename(file);
    expect(result.name).toBe(` -${expectedTimestamp}.pdf`);
    expect(result.type).toBe('image/pdf');
    expect(result instanceof File).toBe(true);
  });
  it('empty filename with extension', () => {
    const file = new File([''], '.pdf', { type: 'image/pdf' });
    const result = appendTimestampToFilename(file);
    expect(result.name).toBe(`-${expectedTimestamp}.pdf`);
    expect(result.type).toBe('image/pdf');
    expect(result instanceof File).toBe(true);
  });
  it('empty filename and no extension', () => {
    const file = new File([''], '', { type: 'image/pdf' });
    const result = appendTimestampToFilename(file);
    expect(result.name).toBe(`-${expectedTimestamp}`);
    expect(result.type).toBe('image/pdf');
    expect(result instanceof File).toBe(true);
  });
  it('filename no extension', () => {
    const file = new File([''], 'helloworld', { type: 'image/pdf' });
    const result = appendTimestampToFilename(file);
    expect(result.name).toBe(`helloworld-${expectedTimestamp}`);
    expect(result.type).toBe('image/pdf');
    expect(result instanceof File).toBe(true);
  });
  // Test each file type in fileTypeMap
  Object.entries(fileTypeMap).forEach(([mimeType, extension]) => {
    it(`appends timestamp to filename with ${extension} extension`, () => {
      const file = new File([''], `test-file.${extension}`, { type: mimeType });
      const result = appendTimestampToFilename(file);
      expect(result.name).toBe(`test-file-${expectedTimestamp}.${extension}`);
      expect(result.type).toBe(mimeType);
      expect(result instanceof File).toBe(true);
    });
  });
});
