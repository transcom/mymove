import React from 'react';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import AsyncPacketDownloadLink, { onPacketDownloadSuccessHandler } from './AsyncPacketDownloadLink';
import { setShowLoadingSpinner } from 'store/general/actions';
import { renderWithProviders } from 'testUtils';

jest.mock('store/general/actions', () => ({
  ...jest.requireActual('store/general/actions'),
  setShowLoadingSpinner: jest.fn().mockImplementation(() => ({
    type: '',
    showSpinner: false,
    loadingSpinnerMessage: '',
  })),
}));

describe('AsyncPacketDownloadLink success', () => {
  it('success', async () => {
    const asyncRetrieval = jest.fn().mockImplementation(() => Promise.resolve());
    const onSuccessHandler = jest.fn();
    const onErrorHandler = jest.fn();
    const expectedId = 'testID';
    const expectedLabel = 'test';

    renderWithProviders(
      <AsyncPacketDownloadLink
        id={expectedId}
        label={expectedLabel}
        asyncRetrieval={asyncRetrieval}
        onSuccess={onSuccessHandler}
        onFailure={onErrorHandler}
      />,
    );
    expect(screen.getByText(expectedLabel, { exact: false })).toBeInTheDocument();

    const downloadButton = screen.getByText(expectedLabel);
    expect(downloadButton).toBeInTheDocument();
    await userEvent.click(downloadButton);

    await waitFor(() => {
      expect(asyncRetrieval).toHaveBeenCalledTimes(1);
      expect(onSuccessHandler).toHaveBeenCalledTimes(1);
      expect(onErrorHandler).toHaveBeenCalledTimes(0);
      expect(setShowLoadingSpinner).toHaveBeenCalled();
    });
  });

  it('AsyncPacketDownloadLink failure', async () => {
    const asyncRetrieval = jest.fn().mockImplementation(() => Promise.reject());
    const onSuccessHandler = jest.fn();
    const onErrorHandler = jest.fn();
    const expectedId = 'testID';
    const expectedLabel = 'test';
    renderWithProviders(
      <AsyncPacketDownloadLink
        id={expectedId}
        label={expectedLabel}
        asyncRetrieval={asyncRetrieval}
        onSucccess={onSuccessHandler}
        onFailure={onErrorHandler}
      />,
    );
    expect(screen.getByText(expectedLabel, { exact: false })).toBeInTheDocument();

    const downloadButton = screen.getByText(expectedLabel);
    expect(downloadButton).toBeInTheDocument();
    await userEvent.click(downloadButton);

    await waitFor(() => {
      expect(asyncRetrieval).toHaveBeenCalledTimes(1);
      expect(asyncRetrieval).toHaveBeenCalledWith(expectedId);
      expect(onSuccessHandler).toHaveBeenCalledTimes(0);
      expect(onErrorHandler).toHaveBeenCalledTimes(1);
      expect(setShowLoadingSpinner).toHaveBeenCalled();
    });
  });

  it('success - downloadPacketOnSuccessHandler', () => {
    const expectedResponseData = 'MOCK_PDF_DATA';
    const expectedFileName = 'test.pdf';
    const expectedContentType = 'application/pdf';

    global.URL.createObjectURL = jest.fn();

    const mockResponse = {
      ok: true,
      headers: {
        'content-disposition': `filename="${expectedFileName}"`,
        'content-type': expectedContentType,
      },
      status: 200,
      data: expectedResponseData,
    };

    function makeAnchor(target) {
      const setAttributeMock = jest.fn((key, value) => (target[key] = value));
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

    onPacketDownloadSuccessHandler(mockResponse);

    // verify response.data is used for blob
    expect(blobSpy).toBeCalledWith([expectedResponseData], {
      type: expectedContentType,
    });

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
