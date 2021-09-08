/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen } from '@testing-library/react';

import DocsUploaded from '.';

describe('DocsUploaded component', () => {
  it.each([
    [1, [{ id: '1', filename: 'The fellowship of the file' }], '1 File uploaded'],
    [
      2,
      [
        { id: '1', filename: 'The twin files' },
        { id: '2', filename: 'The return of the file' },
      ],
      '2 Files uploaded',
    ],
  ])(
    'renders document list with expected heading and number of files (%s)',
    async (numFiles, files, expectedHeading) => {
      render(<DocsUploaded files={files} />);

      const docCountHeading = await screen.findByRole('heading', { level: 6, name: expectedHeading });
      expect(docCountHeading).toBeInTheDocument();

      expect(screen.getAllByTestId('doc-list-item').length).toBe(numFiles);
    },
  );
});
