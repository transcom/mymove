/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { act } from 'react-dom/test-utils';

import FileUpload from './FileUpload';

import { UPLOAD_SCAN_STATUS } from 'shared/constants';
import { waitForAvScan } from 'services/internalApi';

const mockSetOptions = jest.fn();

jest.mock('react-filepond', () => {
  // Mock is before imports, give scope to react or it'll err
  // eslint-disable-next-line global-require
  const MockReact = require('react');

  const setupFilePondMock = () => {
    const root = global.document.createElement('div');

    const item = global.document.createElement('div');
    item.className = 'filepond--item';

    const name = global.document.createElement('span');
    name.className = 'filepond--file-info-main';
    name.textContent = 'dummy.jpg';

    const statusMain = global.document.createElement('span');
    statusMain.className = 'filepond--file-status-main';
    statusMain.textContent = 'Uploading';

    const statusSub = global.document.createElement('span');
    statusSub.className = 'filepond--file-status-sub';
    statusSub.textContent = 'Tap to abort';

    item.append(name, statusMain, statusSub);
    root.appendChild(item);
    return { root, statusMain, statusSub };
  };

  const FilePond = MockReact.forwardRef((props, ref) => {
    const fake = MockReact.useMemo(setupFilePondMock, []);

    // Filepond API
    MockReact.useImperativeHandle(ref, () => ({
      _pond: {
        element: fake.root,
        setOptions: mockSetOptions,
      },
      setOptions: mockSetOptions,
    }));

    return <div data-testid="filepond-stub" ref={(node) => node && node.appendChild(fake.root)} />;
  });

  FilePond.displayName = 'FilePond';

  return { FilePond, registerPlugin: jest.fn() };
});

jest.mock('services/internalApi', () => {
  return {
    waitForAvScan: jest.fn(),
    deleteUpload: jest.fn(),
  };
});

const flushPromises = () =>
  new Promise((resolve) => {
    setTimeout(resolve, 0);
  });

const returnNewDummyFile = () => new File(['123'], 'dummy.jpg', { type: 'image/jpeg', lastModified: Date.now() });

describe('FileUpload processing', () => {
  afterEach(() => {
    jest.clearAllMocks();
  });

  it('changes status label from Uploading -> Scanning for viruses', async () => {
    const createUpload = jest.fn(() => Promise.resolve({ id: 'abc123' }));
    waitForAvScan.mockResolvedValueOnce({ status: 'CLEAN' });

    const wrapper = mount(<FileUpload createUpload={createUpload} />);

    // initial state should read "Uploading"
    const statusBefore = wrapper
      .find('[data-testid="filepond-stub"]')
      .getDOMNode()
      .querySelector('.filepond--file-status-main').textContent;
    expect(statusBefore).toBe('Uploading');

    const { process } = wrapper.find('FilePond').prop('server');

    const loadFunc = jest.fn();
    const errorFunc = jest.fn();

    await act(async () => {
      process('file', returnNewDummyFile(), {}, loadFunc, errorFunc, jest.fn(), jest.fn());
      await flushPromises();
    });

    wrapper.update();

    // grab the new status from stubbed filepond
    const statusAfter = wrapper
      .find('[data-testid="filepond-stub"]')
      .getDOMNode()
      .querySelector('.filepond--file-status-main').textContent;

    // assertions
    expect(statusAfter).toBe('Scanning');
    expect(createUpload).toHaveBeenCalledTimes(1);
    expect(waitForAvScan).toHaveBeenCalledTimes(1);
    expect(loadFunc).toHaveBeenCalledWith('abc123');
    expect(errorFunc).not.toHaveBeenCalled();
  });

  it('it shows file failure when av scan returns INFECTED', async () => {
    const createUpload = jest.fn(() => Promise.resolve({ id: 'abc123' }));
    waitForAvScan.mockRejectedValueOnce(new Error(UPLOAD_SCAN_STATUS.LEGACY_INFECTED));

    const wrapper = mount(<FileUpload createUpload={createUpload} />);

    const { process } = wrapper.find('FilePond').prop('server');
    const errorFunc = jest.fn();

    await act(async () => {
      process('file', returnNewDummyFile(), {}, jest.fn(), errorFunc, jest.fn(), jest.fn());
      await flushPromises();
    });

    expect(errorFunc).toHaveBeenCalledWith('File failed virus scan');
    expect(mockSetOptions).toHaveBeenCalledWith({
      labelFileProcessing: 'File failed virus scan',
    });
  });
});
