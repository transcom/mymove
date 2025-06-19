/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

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
    rotation: 0,
  },
  {
    id: 2,
    filename: 'Test File 2.jpg',
    contentType: 'image/jpeg',
    url: sampleJPG,
    createdAt: '2021-06-12T15:09:26.979879Z',
    rotation: 3,
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

const mockSetRotationValue = jest.fn();
const mockSaveRotation = jest.fn();

jest.mock('@transcom/react-file-viewer', () => ({
  __esModule: true,
  default: ({ fileType, filePath, rotationValue, renderControls }) => {
    const rotateLeft = () => mockSetRotationValue((rotationValue + 270) % 360);
    const rotateRight = () => mockSetRotationValue((rotationValue + 90) % 360);

    return (
      <div>
        <div data-testid="fileTypeText">{fileType}</div>
        <div data-testid="filePathText">{filePath}</div>
        <div data-testid="rotationText">{rotationValue}</div>
        {renderControls &&
          renderControls({
            handleZoomIn: jest.fn(),
            handleZoomOut: jest.fn(),
            handleRotateLeft: () => mockSetRotationValue(rotateLeft()),
            handleRotateRight: () => mockSetRotationValue(rotateRight()),
          })}
        <button type="button" data-testid="rotateLeftButton" onClick={rotateLeft}>
          Rotate left
        </button>
        <button type="button" data-testid="rotateRightButton" onClick={rotateRight}>
          Rotate right
        </button>
      </div>
    );
  },
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

  it('renders the FileViewer with rotation value prop for non PDF file types', () => {
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

  it('renders with rotation value for PDFs', () => {
    const mockPDF = mockFiles[0];
    render(<DocViewerContent fileType={mockPDF.contentType} filePath={mockPDF.url} rotationValue={mockPDF.rotation} />);

    expect(screen.getByTestId('fileTypeText').textContent).toContain('application/pdf');
    expect(screen.getByTestId('filePathText').textContent).toContain(samplePDF);
    expect(screen.getByTestId('rotationText').textContent).toContain('0');
  });

  it('calls setRotationValue when rotate left button is clicked for PDFs', async () => {
    const mockPDF = mockFiles[0];
    const startingRotation = 90;
    render(
      <DocViewerContent
        fileType={mockPDF.contentType}
        filePath={mockPDF.url}
        rotationValue={startingRotation}
        setRotationValue={mockSetRotationValue}
        saveRotation={mockSaveRotation}
      />,
    );
    const rotateLeftBtn = screen.getByTestId('rotateLeftButton');
    await userEvent.click(rotateLeftBtn);
    expect(mockSetRotationValue).toHaveBeenCalledWith(0);
  });

  it('calls setRotationValue when rotate right button is clicked for PDFs', async () => {
    const mockPDF = mockFiles[0];
    const startingRotation = 90;
    render(
      <DocViewerContent
        fileType={mockPDF.contentType}
        filePath={mockPDF.url}
        rotationValue={startingRotation}
        setRotationValue={mockSetRotationValue}
        saveRotation={mockSaveRotation}
      />,
    );
    const rotateRightBtn = screen.getByTestId('rotateRightButton');
    await userEvent.click(rotateRightBtn);
    expect(mockSetRotationValue).toHaveBeenCalledWith(180);
  });

  it('calls saveRotation when Save is clicked and button is enabled', async () => {
    render(
      <DocViewerContent
        fileType="pdf"
        filePath="/test.pdf"
        rotationValue={90}
        disableSaveButton={false}
        setRotationValue={mockSetRotationValue}
        saveRotation={mockSaveRotation}
      />,
    );
    const saveBtn = screen.getByRole('button', { name: /save/i });
    await userEvent.click(saveBtn);
    expect(mockSaveRotation).toHaveBeenCalled();
  });

  it('disables the Save button when disableSaveButton is true', () => {
    render(
      <DocViewerContent
        fileType="pdf"
        filePath="/test.pdf"
        rotationValue={90}
        disableSaveButton
        setRotationValue={mockSetRotationValue}
        saveRotation={mockSaveRotation}
      />,
    );
    const saveBtn = screen.getByRole('button', { name: /save/i });
    expect(saveBtn).toBeDisabled();
  });
});
