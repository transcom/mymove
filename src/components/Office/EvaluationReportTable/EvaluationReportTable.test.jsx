import React from 'react';
import { render, screen } from '@testing-library/react';

import EvaluationReportTable from './EvaluationReportTable';

import { MockProviders } from 'testUtils';

const submittedReport = {
  id: 'a7fdb0b3-f876-4686-b94f-ad20a2c9a63d',
  location: 'DESTINATION',
  moveID: 'bd1bbbdc-1710-4831-aa41-e35ebedff5cd',
  shipmentID: '38e87840-d385-413e-9746-b2da2c8245bb',
  submittedAt: '2022-07-14T19:21:27.565Z',
  type: 'SHIPMENT',
  violationsObserved: true,
};
const draftReport = {
  id: '1f9d18a8-7688-428d-be8e-3f3c59a0ae5e',
  location: null,
  moveID: 'bd1bbbdc-1710-4831-aa41-e35ebedff5cd',
  shipmentID: 'd46825dd-cf90-442b-96de-c5075ea2f1bf',
  submittedAt: null,
  type: 'SHIPMENT',
  violationsObserved: true,
};

const customerInfo = {
  agency: 'ARMY',
  backup_contact: { email: 'email@example.com', name: 'name', phone: '555-555-5555' },
  current_address: {
    city: 'Beverly Hills',
    country: 'US',
    eTag: 'MjAyMi0wOC0xNVQxNjoxMToyNi4zMzIwOTFa',
    id: '28f11990-7ced-4d01-87ad-b16f2c86ea83',
    postalCode: '90210',
    state: 'CA',
    streetAddress1: '123 Any Street',
    streetAddress2: 'P.O. Box 12345',
    streetAddress3: 'c/o Some Person',
  },
  dodID: '5052247544',
  eTag: 'MjAyMi0wOC0xNVQxNjoxMToyNi4zNTkzNFo=',
  email: 'leo_spaceman_sm@example.com',
  first_name: 'Leo',
  id: 'ea557b1f-2660-4d6b-89a0-fb1b5efd2113',
  last_name: 'Spacemen',
  phone: '555-555-5555',
  userID: 'f4bbfcdf-ef66-4ce7-92f8-4c1bf507d596',
};

describe('EvaluationReportTable', () => {
  it('renders empty table', () => {
    render(
      <MockProviders>
        <EvaluationReportTable
          reports={[]}
          emptyText="no reports for this"
          shipments={[]}
          moveCode="FAKEIT"
          grade="E_4"
          customerInfo={customerInfo}
          deleteReport={jest.fn()}
          setReportToDelete={jest.fn()}
          setIsDeleteModalOpen={jest.fn()}
          isDeleteModalOpen={false}
        />
      </MockProviders>,
    );
    expect(screen.getByText('Report ID')).toBeInTheDocument();
    expect(screen.getByText('Date submitted')).toBeInTheDocument();
    expect(screen.getByText('Location')).toBeInTheDocument();
    expect(screen.getByText('Violations')).toBeInTheDocument();
    expect(screen.getByText('Serious Incident')).toBeInTheDocument();
    expect(screen.getByText('no reports for this')).toBeInTheDocument();
  });
  it('renders table with a draft report', () => {
    render(
      <MockProviders>
        <EvaluationReportTable
          reports={[draftReport]}
          emptyText="no reports for this"
          shipments={[]}
          moveCode="FAKEIT"
          grade="E_4"
          customerInfo={customerInfo}
          deleteReport={jest.fn()}
          setReportToDelete={jest.fn()}
          setIsDeleteModalOpen={jest.fn()}
          isDeleteModalOpen={false}
        />
      </MockProviders>,
    );
    expect(screen.getByText('Report ID')).toBeInTheDocument();
    expect(screen.getByText('Date submitted')).toBeInTheDocument();
    expect(screen.getByText('Location')).toBeInTheDocument();
    expect(screen.getByText('Violations')).toBeInTheDocument();
    expect(screen.getByText('Serious Incident')).toBeInTheDocument();
    expect(screen.queryByText('no reports for this')).not.toBeInTheDocument();

    expect(screen.getByTestId('tag')).toHaveTextContent('DRAFT');

    expect(screen.getByText('#QA-1F9D1')).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Edit report' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Delete' })).toBeInTheDocument();
  });
  it('renders table with a submitted report', () => {
    render(
      <MockProviders>
        <EvaluationReportTable
          reports={[submittedReport]}
          emptyText="no reports for this"
          shipments={[]}
          moveCode="FAKEIT"
          grade="E_4"
          customerInfo={customerInfo}
          deleteReport={jest.fn()}
          setReportToDelete={jest.fn()}
          setIsDeleteModalOpen={jest.fn()}
          isDeleteModalOpen={false}
        />
      </MockProviders>,
    );
    expect(screen.getByText('Report ID')).toBeInTheDocument();
    expect(screen.getByText('Date submitted')).toBeInTheDocument();
    expect(screen.getByText('Location')).toBeInTheDocument();
    expect(screen.getByText('Violations')).toBeInTheDocument();
    expect(screen.getByText('Serious Incident')).toBeInTheDocument();

    expect(screen.queryByText('no reports for this')).not.toBeInTheDocument();
    expect(screen.queryByTestId('tag')).not.toBeInTheDocument();

    expect(screen.getByText('#QA-A7FDB')).toBeInTheDocument();
    expect(screen.getByText('14 Jul 2022')).toBeInTheDocument();
    expect(screen.getByText('Destination')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'View report' })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Download' })).toBeInTheDocument();
  });
});
