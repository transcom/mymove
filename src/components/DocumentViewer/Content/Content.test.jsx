/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen } from '@testing-library/react';

import samplePDF from '../sample.pdf';
import sampleJPG from '../sample.jpg';
import samplePNG from '../sample2.png';
import sampleGIF from '../sample3.gif';

import DocViewerContent from './Content';

const mockFiles = [
  {
    id: 1,
    filename: 'Test File.pdf',
    contentType: 'application/pdf',
    url: samplePDF,
    createdAt: '2021-06-14T15:09:26.979879Z',
  },
  {
    id: 2,
    filename: 'Test File 2.jpg',
    contentType: 'image/jpeg',
    url: sampleJPG,
    createdAt: '2021-06-12T15:09:26.979879Z',
  },
  {
    id: 3,
    filename: 'Test File 3.png',
    contentType: 'image/png',
    url: samplePNG,
    createdAt: '2021-06-15T15:09:26.979879Z',
    rotation: 1,
  },
  {
    id: 4,
    filename: 'Test File 4.gif',
    contentType: 'image/gif',
    url: sampleGIF,
    createdAt: '2021-06-16T15:09:26.979879Z',
    rotation: 3,
  },
];

jest.mock('@transcom/react-file-viewer', () => ({
  __esModule: true,
  default: ({ fileType, filePath, rotationValue }) => (
    <div>
      <div>
        <div data-testid="fileTypeText">{fileType}</div>
        <div data-testid="filePathText">{filePath}</div>
        <div data-testid="rotationText">{rotationValue}</div>
      </div>
    </div>
  ),
}));

describe('DocViewerContent', () => {
  it('renders without crashing', () => {
    render(<DocViewerContent />);

    expect(screen.getByTestId('DocViewerContent')).toBeInTheDocument();
  });

  it('renders the FileViewer with the file props', () => {
    render(<DocViewerContent fileType={mockFiles[0].contentType} filePath={mockFiles[0].url} />);

    expect(screen.getByTestId('fileTypeText').textContent).toContain('application/pdf');
    expect(screen.getByTestId('filePathText').textContent).toContain(samplePDF);
  });

  it('renders the FileViewer with rotation value prop', () => {
    render(
      <DocViewerContent
        fileType={mockFiles[2].contentType}
        filePath={mockFiles[2].url}
        rotationValue={mockFiles[2].rotation}
      />,
    );

    expect(screen.getByTestId('fileTypeText').textContent).toContain('image/png');
    expect(screen.getByTestId('filePathText').textContent).toContain(samplePNG);
    expect(screen.getByTestId('rotationText').textContent).toContain('1');
  });
});
