import { downloadPPMAOAPacketOnSuccessHandler } from './download';

afterEach(() => {
  jest.resetAllMocks();
});

describe('downloadPPMAOAPacketOnSuccessHandler', () => {
  it('success - downloadPPMAOAPacketOnSuccessHandler', () => {
    const expectedResponseData = 'MOCK_PDF_DATA';
    const expectedFileName = 'test.pdf';

    global.URL.createObjectURL = jest.fn();

    const mockResponse = {
      ok: true,
      headers: {
        'content-disposition': `filename="${expectedFileName}"`,
      },
      status: 200,
      data: expectedResponseData,
    };

    function makeAnchor(target) {
      /* eslint-disable no-param-reassign, no-return-assign */
      const setAttributeMock = jest.fn((key, value) => (target[key] = value));
      /* eslint-enable no-param-reassign, no-return-assign */
      return {
        target,
        setAttribute: setAttributeMock,
        click: jest.fn(),
        remove: jest.fn(),
        parentNode: {
          removeChild: jest.fn(),
        },
      };
    }

    jest.spyOn(document.body, 'appendChild').mockReturnValue(null);
    const anchor = makeAnchor({ href: '#', download: '' });
    jest.spyOn(document, 'createElement').mockReturnValue(anchor);
    const clickSpy = jest.spyOn(anchor, 'click');

    const mBlob = { size: 1024, type: 'application/pdf' };
    const blobSpy = jest.spyOn(global, 'Blob').mockImplementationOnce(() => mBlob);

    downloadPPMAOAPacketOnSuccessHandler(mockResponse);

    // verify response.data is used for blob
    expect(blobSpy).toBeCalledWith([expectedResponseData]);

    // verify hyperlink was created
    expect(document.createElement).toBeCalledWith('a');

    // verify download attribute is from content-disposition
    expect(document.body.appendChild).toBeCalledWith(
      expect.objectContaining({
        target: { download: expectedFileName, href: '#' },
      }),
    );

    // verify click event is invoked to download file
    expect(clickSpy).toHaveBeenCalledTimes(1);

    // verify link is removed
    expect(anchor.parentNode.removeChild).toBeCalledWith(anchor);
  });
});
