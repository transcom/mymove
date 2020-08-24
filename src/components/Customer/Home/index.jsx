/* eslint-disable no-console */
import React from 'react';
import { Alert } from '@trussworks/react-uswds';

import Helper from './Helper';
import Step from './Step';
import DocsUploaded from './DocsUploaded';
import ShipmentList from './ShipmentList';
import Footer from './Footer';
import styles from './Home.module.scss';

import { SHIPMENT_OPTIONS } from 'shared/constants';

const shipments = [
  { type: SHIPMENT_OPTIONS.PPM, id: '#123ABC-001' },
  { type: SHIPMENT_OPTIONS.HHG, id: '#123ABC-002' },
  { type: SHIPMENT_OPTIONS.NTS, id: '#123ABC-003' },
];

const Home = () => {
  function handleShipmentClick(shipment) {
    console.log('this is the shipment', shipment);
  }

  return (
    <div className={`usa-prose grid-container ${styles['grid-container']}`}>
      <header className={`${styles['customer-header']} padding-top-3 padding-bottom-3 margin-bottom-2`}>
        <h2>Riley Baker</h2>
        <p>
          You&apos;re leaving <strong>Buckley AFB</strong>
        </p>
      </header>
      <Alert className="margin-top-2 margin-bottom-2" slim type="success">
        Thank you for adding your Profile information
      </Alert>

      <Helper
        title="Next step: Add your orders"
        helpList={[
          'If you have a hard copy, you can take photos of each page',
          'If you have a PDF, you can upload that',
        ]}
      />
      <Step
        complete
        completedHeaderText="Profile complete"
        description="Make sure to keep your personal information up to date during your move"
        editBtnDisabled
        editBtnLabel="Edit"
        headerText="Profile complete"
        onActionBtnClick={() => console.log('some action')}
        step="1"
        onEditClick={(e) => {
          e.preventDefault();
          console.log('edit clicked');
        }}
      />

      <Step
        complete
        completedHeaderText="Orders uploaded"
        description="Upload photos of each page, or upload a PDF"
        editBtnLabel="Edit"
        onEditClick={() => console.log('edit button clicked')}
        headerText="Upload orders"
        onActionBtnClick={() => console.log('some action')}
        step="2"
      >
        <DocsUploaded
          files={[
            { filename: 'Screen Shot 2020-09-11 at 12.56.58 PM.png' },
            { filename: 'Screen Shot 2020-09-11 at 12.58.12 PM.png' },
            { filename: 'orderspage3_20200723.png' },
          ]}
        />
      </Step>

      <Step
        actionBtnLabel="Add another shipment"
        complete
        completedHeaderText="Shipments"
        description="Tell us where you're going and when you want to get there. We'll help you set up shipments to make it work"
        headerText="Shipments"
        secondaryBtn
        step="3"
      >
        <ShipmentList shipments={shipments} onShipmentClick={handleShipmentClick} />
      </Step>

      <Step
        actionBtnDisabled
        actionBtnLabel="Review and submit"
        containerClassName="margin-bottom-8"
        description="Review your move details and sign the legal paperwork, then send the info on to your move counselor"
        headerText="Confirm move request"
        onActionBtnClick={() => console.log('some action')}
        step="4"
      />
      <Footer
        header="Contacts"
        dutyStationName="Seymour Johnson AFB"
        officeType="Origin Transportation Office"
        telephone="(919) 722-5458"
      />
    </div>
  );
};

export default Home;
