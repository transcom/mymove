import React from 'react';
import { render, screen } from '@testing-library/react';
import { Formik } from 'formik';

import ShipmentFormRemarks from './ShipmentFormRemarks';

import { roleTypes } from 'constants/userRoles';
import { SHIPMENT_OPTIONS } from 'shared/constants';

describe('components/Office/ShipmentFormRemarks', () => {
  it('renders correctly as a Service Counselor', () => {
    render(
      <Formik initialValues={{ counselorRemarks: 'Counselor remarks from initial values' }}>
        <ShipmentFormRemarks
          userRole={roleTypes.SERVICES_COUNSELOR}
          shipmentType={SHIPMENT_OPTIONS.HHG}
          customerRemarks="Customer remarks from props"
          counselorRemarks="Counselor remarks from props"
        />
      </Formik>,
    );

    expect(screen.getByRole('textbox')).toBeInTheDocument();
    expect(screen.getByText('Customer remarks from props')).toBeInTheDocument();
    expect(screen.getByText('Counselor remarks from initial values')).toBeInTheDocument();
  });

  it('renders correctly as a Service Counselor with a PPM Shipment', () => {
    render(
      <Formik initialValues={{ counselorRemarks: 'Counselor remarks from initial values' }}>
        <ShipmentFormRemarks
          userRole={roleTypes.SERVICES_COUNSELOR}
          shipmentType={SHIPMENT_OPTIONS.PPM}
          customerRemarks="Customer remarks from props"
          counselorRemarks="Counselor remarks from props"
        />
      </Formik>,
    );

    expect(screen.getByRole('textbox')).toBeInTheDocument();
    expect(screen.queryByText('Customer remarks from props')).not.toBeInTheDocument();
    expect(screen.getByText('Counselor remarks from initial values')).toBeInTheDocument();
  });

  it('renders required asterisk for Counselor remarks as SC rejecting advance on PPM', () => {
    render(
      <Formik initialValues={{}}>
        <ShipmentFormRemarks
          userRole={roleTypes.SERVICES_COUNSELOR}
          shipmentType={SHIPMENT_OPTIONS.PPM}
          counselorRemarks="Counselor remarks from props"
          advanceStatus="REJECTED"
          isAdvancePage
        />
      </Formik>,
    );

    expect(document.querySelector('#reqAsteriskMsg')).toHaveTextContent('Fields marked with * are required.');
    expect(screen.getByLabelText('Counselor remarks *')).toBeInTheDocument();
  });

  it('does not render required asterisk for Counselor remarks as SC approving advance on PPM', () => {
    render(
      <Formik initialValues={{}}>
        <ShipmentFormRemarks
          userRole={roleTypes.SERVICES_COUNSELOR}
          shipmentType={SHIPMENT_OPTIONS.PPM}
          counselorRemarks="Counselor remarks from props"
          advanceStatus="APPROVED"
        />
      </Formik>,
    );

    expect(screen.getByLabelText('Counselor remarks')).toBeInTheDocument();
  });

  it('renders correctly as a TOO', () => {
    render(
      <Formik initialValues={{}}>
        <ShipmentFormRemarks
          userRole={roleTypes.TOO}
          shipmentType={SHIPMENT_OPTIONS.HHG}
          counselorRemarks="Counselor remarks from props"
        />
      </Formik>,
    );

    expect(screen.queryByRole('textbox')).not.toBeInTheDocument();
    expect(screen.getByText('—')).toBeInTheDocument();
    expect(screen.getByText('Counselor remarks from props')).toBeInTheDocument();
  });
});
