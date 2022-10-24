import { React } from 'react';
import { fireEvent, render, screen } from '@testing-library/react';

import UploadsTable from './UploadsTable';

describe('UploadTable component', () => {
  const testPropsNoUploads = {
    onDelete: jest.fn(),
    uploads: [],
  };

  const testPropsMultipleUploads = {
    onDelete: jest.fn(),
    uploads: [
      {
        bytes: 9043,
        contentType: 'image/png',
        createdAt: '2021-06-21T19:51:49.441Z',
        filename: 'orders1.png',
        id: '00000000-0000-0000-0000-000000000001',
        url: '',
      },
      {
        bytes: 4043,
        contentType: 'application/pdf',
        createdAt: '2021-06-21T20:33:22.724Z',
        filename: 'orders2.pdf',
        id: '00000000-0000-0000-0000-000000000002',
        url: '',
      },
    ],
  };

  it('renders nothing if there are no uploads', async () => {
    render(<UploadsTable {...testPropsNoUploads} />);
    expect(screen.queryByRole('heading', { name: /0 FILES UPLOADED/i, level: 6 })).toBeNull();
    expect(screen.queryAllByRole('listitem')).toHaveLength(0);
  });

  it('renders if there are uploads', async () => {
    render(<UploadsTable {...testPropsMultipleUploads} />);
    expect(screen.getByRole('heading', { name: /2 FILES UPLOADED/i, level: 6 })).toBeInTheDocument();

    expect(screen.getAllByRole('listitem')).toHaveLength(testPropsMultipleUploads.uploads.length);
    expect(screen.getAllByRole('button', { name: 'Delete' })).toHaveLength(testPropsMultipleUploads.uploads.length);
  });

  it('Delete button calls onDelete callback', async () => {
    render(<UploadsTable {...testPropsMultipleUploads} />);
    const deleteButtons = screen.getAllByRole('button', { name: 'Delete' });
    fireEvent.click(deleteButtons[0]);
    expect(testPropsMultipleUploads.onDelete).toHaveBeenCalled();
  });
});
