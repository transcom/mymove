import { React } from 'react';
import { screen, waitFor } from '@testing-library/react';

import AdditionalDocuments from './AdditionalDocuments';

import { renderWithProviders } from 'testUtils';
import { getMove } from 'services/internalApi';
import { selectCurrentMove } from 'store/entities/selectors';

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectCurrentMove: jest.fn(),
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  getMove: jest.fn().mockImplementation(() => Promise.resolve()),
}));

const testMove = {
  additionalDocuments: {
    id: 'c43ae36e-4e15-4cb3-865a-e4dccffa0df7',
    service_member_id: 'dfdd3e21-3988-4104-a5c2-06b195f9b7f0',
    uploads: [
      {
        bytes: 120653,
        contentType: 'application/pdf',
        createdAt: '2024-05-29T19:14:39.108Z',
        filename: '9380-Statement-20240430.pdf',
        id: 'c3c0cda9-a77e-4b8b-8b8b-67ccadc3c862',
        status: 'PROCESSING',
        updatedAt: '2024-05-29T19:14:39.108Z',
        url: '/storage/user/accf760b-2e3d-4af8-a59b-c10b591dcc15/uploads/c3c0cda9-a77e-4b8b-8b8b-67ccadc3c862?contentType=application%2Fpdf',
      },
      {
        bytes: 307051,
        contentType: 'image/png',
        createdAt: '2024-05-30T04:23:27.241Z',
        filename: 'Screenshot 2024-05-16 at 3.33.52 PM.png',
        id: '70a35ab0-a3f5-44a3-8702-0bb7d0c568c8',
        status: 'PROCESSING',
        updatedAt: '2024-05-30T04:23:27.241Z',
        url: '/storage/user/accf760b-2e3d-4af8-a59b-c10b591dcc15/uploads/70a35ab0-a3f5-44a3-8702-0bb7d0c568c8?contentType=image%2Fpng',
      },
      {
        bytes: 82301,
        contentType: 'image/png',
        createdAt: '2024-05-30T04:33:10.622Z',
        filename: 'Screenshot 2024-05-17 at 1.09.21 PM.png',
        id: 'b11c0130-2403-4287-b464-4c5ac17797b3',
        status: 'PROCESSING',
        updatedAt: '2024-05-30T04:33:10.622Z',
        url: '/storage/user/accf760b-2e3d-4af8-a59b-c10b591dcc15/uploads/b11c0130-2403-4287-b464-4c5ac17797b3?contentType=image%2Fpng',
      },
    ],
  },
  created_at: '2024-05-29T18:46:17.808Z',
  eTag: 'MjAyNC0wNS0yOVQxOToxNDozOS4xMDQyNzJa',
  id: '43a369e8-5fa3-4a13-9d9a-36d86731c1da',
  locator: '988HDJ',
  mto_shipments: ['c93bf4d1-1470-4c50-b2b6-f736abd2986a'],
  orders_id: '69967de3-3d9d-4e73-a497-f401884393bf',
  primeCounselingCompletedAt: '0001-01-01T00:00:00.000Z',
  service_member_id: 'dfdd3e21-3988-4104-a5c2-06b195f9b7f0',
  status: 'NEEDS SERVICE COUNSELING',
  submitted_at: '2024-05-29T18:47:26.360Z',
  updated_at: '2024-05-29T19:14:39.104Z',
};

describe('Additional Documents component', () => {
  const testProps = {
    move: testMove,
    updateMove: jest.fn(),
  };
  it('renders all content of AdditionalDocuments', () => {
    selectCurrentMove.mockImplementation(() => testMove);
    getMove.mockResolvedValue(testMove);
    renderWithProviders(<AdditionalDocuments {...testProps} />);

    expect(screen.getByLabelText('Upload')).toBeInTheDocument();
  });
});
